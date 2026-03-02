package diagon

/*
#include <stdlib.h>
#include <string.h>
#include "diagon/c_api/diagon_c_api.h"
*/
import "C"

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"unsafe"
)

// IndexingTimingResult holds timing breakdown from BulkIndexWithTiming.
type IndexingTimingResult struct {
	PrepTime time.Duration // Go→C conversion (flattenMap + CString + field creation)
	AddTime  time.Duration // Pure C++ diagon_add_documents
	FreeTime time.Duration // Document handle cleanup
	Total    time.Duration // PrepTime + AddTime + FreeTime
	NumDocs  int
}

// BulkIndexWithTiming indexes documents like BulkIndexDocuments but returns per-phase timing.
func (s *Shard) BulkIndexWithTiming(docs []struct {
	ID         string
	Doc        map[string]interface{}
	SourceJSON []byte
}) (*IndexingTimingResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	numDocs := len(docs)
	if numDocs == 0 {
		return &IndexingTimingResult{}, nil
	}

	initCachedFieldNames()

	result := &IndexingTimingResult{NumDocs: numDocs}

	cDocsRaw := C.malloc(C.size_t(numDocs) * C.size_t(unsafe.Sizeof(uintptr(0))))
	defer C.free(cDocsRaw)
	cDocs := unsafe.Slice((*C.DiagonDocument)(cDocsRaw), numDocs)

	fieldNameCache := make(map[string]*C.char, 32)
	tmpPtrs := make([]unsafe.Pointer, 0, 20)

	// Phase: Go→C conversion
	prepStart := time.Now()
	for i, item := range docs {
		tmpPtrs = tmpPtrs[:0]
		diagonDoc := C.diagon_create_document()

		cDocID := C.CString(item.ID)
		tmpPtrs = append(tmpPtrs, unsafe.Pointer(cDocID))

		idField := C.diagon_create_string_field(cachedIDFieldName, cDocID)
		C.diagon_document_add_field(diagonDoc, idField)
		storedIDField := C.diagon_create_stored_field(cachedIDFieldName, cDocID)
		C.diagon_document_add_field(diagonDoc, storedIDField)

		sourceBytes := item.SourceJSON
		if sourceBytes == nil {
			sourceBytes, _ = json.Marshal(item.Doc)
		}
		if len(sourceBytes) > 0 {
			cSourceValue := cStringFromBytes(sourceBytes)
			tmpPtrs = append(tmpPtrs, unsafe.Pointer(cSourceValue))
			sourceField := C.diagon_create_stored_field(cachedSourceFieldName, cSourceValue)
			C.diagon_document_add_field(diagonDoc, sourceField)
		}

		flatDoc := flattenMap("", item.Doc)
		for key, value := range flatDoc {
			cFieldName, ok := fieldNameCache[key]
			if !ok {
				cFieldName = C.CString(key)
				fieldNameCache[key] = cFieldName
			}

			switch v := value.(type) {
			case string:
				if epochMs, isDate := tryParseDateToEpochMs(v); isDate {
					doubleField := C.diagon_create_indexed_double_field(cFieldName, C.double(float64(epochMs)))
					C.diagon_document_add_field(diagonDoc, doubleField)
					cValue := C.CString(v)
					tmpPtrs = append(tmpPtrs, unsafe.Pointer(cValue))
					storedField := C.diagon_create_stored_field(cFieldName, cValue)
					C.diagon_document_add_field(diagonDoc, storedField)
				} else if isKeywordLike(v) {
					cValue := C.CString(v)
					tmpPtrs = append(tmpPtrs, unsafe.Pointer(cValue))
					stringField := C.diagon_create_string_field(cFieldName, cValue)
					C.diagon_document_add_field(diagonDoc, stringField)
					storedField := C.diagon_create_stored_field(cFieldName, cValue)
					C.diagon_document_add_field(diagonDoc, storedField)
				} else {
					cValue := C.CString(v)
					tmpPtrs = append(tmpPtrs, unsafe.Pointer(cValue))
					field := C.diagon_create_text_field(cFieldName, cValue)
					C.diagon_document_add_field(diagonDoc, field)
					stringField := C.diagon_create_string_field(cFieldName, cValue)
					C.diagon_document_add_field(diagonDoc, stringField)
				}
			case int, int32, int64:
				val := int64(0)
				switch iv := v.(type) {
				case int:
					val = int64(iv)
				case int32:
					val = int64(iv)
				case int64:
					val = iv
				}
				field := C.diagon_create_indexed_double_field(cFieldName, C.double(float64(val)))
				C.diagon_document_add_field(diagonDoc, field)
				cValueStr := C.CString(strconv.FormatInt(val, 10))
				tmpPtrs = append(tmpPtrs, unsafe.Pointer(cValueStr))
				storedField := C.diagon_create_stored_field(cFieldName, cValueStr)
				C.diagon_document_add_field(diagonDoc, storedField)
			case float32, float64:
				fval := float64(0)
				switch fv := v.(type) {
				case float32:
					fval = float64(fv)
				case float64:
					fval = fv
				}
				field := C.diagon_create_indexed_double_field(cFieldName, C.double(fval))
				C.diagon_document_add_field(diagonDoc, field)
				cValueStr := C.CString(strconv.FormatFloat(fval, 'f', 6, 64))
				tmpPtrs = append(tmpPtrs, unsafe.Pointer(cValueStr))
				storedField := C.diagon_create_stored_field(cFieldName, cValueStr)
				C.diagon_document_add_field(diagonDoc, storedField)
			default:
				jsonBytes, err := json.Marshal(v)
				if err != nil {
					continue
				}
				cValue := cStringFromBytes(jsonBytes)
				tmpPtrs = append(tmpPtrs, unsafe.Pointer(cValue))
				field := C.diagon_create_stored_field(cFieldName, cValue)
				C.diagon_document_add_field(diagonDoc, field)
			}
		}

		cDocs[i] = diagonDoc
		for _, ptr := range tmpPtrs {
			C.free(ptr)
		}
	}
	result.PrepTime = time.Since(prepStart)

	// Phase: Pure C++ indexing
	addStart := time.Now()
	cResult := C.diagon_add_documents(s.writer, &cDocs[0], C.int(numDocs))
	result.AddTime = time.Since(addStart)

	// Phase: Cleanup
	freeStart := time.Now()
	for i := 0; i < numDocs; i++ {
		C.diagon_free_document(cDocs[i])
	}
	for _, cs := range fieldNameCache {
		C.free(unsafe.Pointer(cs))
	}
	result.FreeTime = time.Since(freeStart)

	result.Total = result.PrepTime + result.AddTime + result.FreeTime

	if int(cResult) < 0 {
		errMsg := C.GoString(C.diagon_last_error())
		return result, fmt.Errorf("batch add_documents failed: %s", errMsg)
	}
	if int(cResult) < numDocs {
		return result, fmt.Errorf("batch add_documents: only %d of %d documents added", int(cResult), numDocs)
	}

	s.readerDirty = true
	return result, nil
}
