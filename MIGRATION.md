# Migration Guide: Quidditch → CONJUGATE

**Date**: January 29, 2026
**Migration Status**: In Progress
**Target Completion**: February 15, 2026

---

## Overview

The project formerly known as "Quidditch" has been rebranded to **CONJUGATE** to avoid trademark conflicts with Warner Bros Entertainment Inc. and to better communicate our technical capabilities.

**New Name**: CONJUGATE - **C**loud-native **O**bservability + **N**atural-language **J**oint **U**nderstanding **G**ranular search **A**nalytics **T**unable **E**ngine

---

## Why the Change?

### Legal Considerations

"Quidditch" is a registered trademark of Warner Bros for the fictional sport from Harry Potter. After legal analysis:
- **70-85% risk** of trademark enforcement
- **Real-world precedent**: The real-world Quidditch sport rebranded to "Quadball" in 2022
- **Open source visibility**: Public projects are more vulnerable to enforcement

### Strategic Benefits

The new name **CONJUGATE**:
- ✅ No trademark conflicts
- ✅ Emphasizes four strategic pillars: Cloud-Native, Observability, NLP, Efficiency
- ✅ Professional identity for enterprise adoption
- ✅ Unique brand in the search engine space

See [NAMING.md](NAMING.md) for complete rationale.

---

## What's Changing

### ✅ Already Updated

- [x] README.md
- [x] Project name and backronym
- [x] NAMING.md documentation
- [x] MIGRATION.md (this file)

### ⏳ In Progress (Week 1-2)

- [ ] All documentation files (ARCHITECTURE, ROADMAP, etc.)
- [ ] Code comments and package descriptions
- [ ] Kubernetes CRD definitions
- [ ] API package paths
- [ ] Command-line tool names
- [ ] Configuration examples
- [ ] Test fixtures

### 📅 Planned (Week 3-4)

- [ ] GitHub repository rename
- [ ] Docker image names
- [ ] Go module path
- [ ] Domain name registration
- [ ] CI/CD pipeline updates
- [ ] Public announcement

---

## For Users: What You Need to Do

### If You Haven't Deployed Yet

**Good news**: Just use the new name! No action required.

```bash
# Use the new repository name
git clone https://github.com/yourorg/conjugate.git
cd conjugate

# Follow the updated README.md
./scripts/deploy-k8s.sh --profile dev
```

---

### If You're Running Quidditch in Development

**Migration Steps**:

#### 1. Update Git Remote (Week 3)

```bash
cd quidditch
git remote set-url origin https://github.com/yourorg/conjugate.git
git pull
```

#### 2. Update Local Configuration

```bash
# Rename your local directory (optional)
cd ..
mv quidditch conjugate
cd conjugate
```

#### 3. Update Kubernetes Resources

```bash
# Update CRD references
# Old:
apiVersion: quidditch.io/v1
kind: QuidditchCluster

# New:
apiVersion: conjugate.io/v1
kind: ConjugateCluster
```

#### 4. Update Application Code

```bash
# Update import paths (when Go module path changes)
# Old:
import "github.com/yourorg/quidditch/pkg/..."

# New:
import "github.com/yourorg/conjugate/pkg/..."
```

#### 5. Update Environment Variables

```bash
# Old:
QUIDDITCH_ENDPOINT=http://localhost:9200

# New:
CONJUGATE_ENDPOINT=http://localhost:9200
```

---

### If You're Running Quidditch in Production

⚠️ **Please wait for stable migration path** ⚠️

We're working on zero-downtime migration strategy. Do NOT migrate production deployments until we publish:

- [ ] Official migration announcement (Week 4)
- [ ] Backward compatibility layer (Week 4)
- [ ] Rolling upgrade procedure (Week 4)
- [ ] Rollback procedures (Week 4)

**Subscribe to updates**: [GitHub Discussions - Migration](https://github.com/yourorg/conjugate/discussions/migration)

---

## Backward Compatibility

### Supported During Transition (Weeks 1-4)

To minimize disruption, we'll support both names during the transition period:

```yaml
# Both CRD names will work
apiVersion: quidditch.io/v1        # ⚠️ Deprecated
kind: QuidditchCluster             # ⚠️ Deprecated

apiVersion: conjugate.io/v1        # ✅ Recommended
kind: ConjugateCluster             # ✅ Recommended
```

```bash
# Both command names will work
quidditch pipeline deploy ...      # ⚠️ Deprecated
conjugate pipeline deploy ...      # ✅ Recommended
```

### Deprecation Timeline

| Date | Status | Action |
|------|--------|--------|
| **Jan 29, 2026** | Rebranding announced | Start using new name |
| **Feb 15, 2026** | Backward compatibility | Both names supported |
| **Mar 15, 2026** | Deprecation warnings | Warnings for old name usage |
| **Apr 15, 2026** | End of support | Old names no longer supported |

---

## Detailed Migration by Component

### 1. Kubernetes Deployments

#### CRD Migration

**Old CRD**:
```yaml
apiVersion: quidditch.io/v1
kind: QuidditchCluster
metadata:
  name: quidditch-prod
  namespace: quidditch
spec:
  version: "1.0.0"
  master:
    replicas: 3
  coordination:
    replicas: 2
  data:
    replicas: 3
```

**New CRD**:
```yaml
apiVersion: conjugate.io/v1
kind: ConjugateCluster
metadata:
  name: conjugate-prod
  namespace: conjugate
spec:
  version: "1.0.0"
  master:
    replicas: 3
  coordination:
    replicas: 2
  data:
    replicas: 3
```

**Migration Command**:
```bash
# 1. Export existing cluster
kubectl get quidditchcluster quidditch-prod -n quidditch -o yaml > old-cluster.yaml

# 2. Transform to new format
sed 's/quidditch/conjugate/g' old-cluster.yaml > new-cluster.yaml

# 3. Apply new cluster (parallel deployment)
kubectl apply -f new-cluster.yaml

# 4. Migrate data (if needed)
# ... data migration steps ...

# 5. Delete old cluster
kubectl delete quidditchcluster quidditch-prod -n quidditch
```

---

### 2. Service Names

**Old Service Names**:
```
quidditch-master
quidditch-coordination
quidditch-data
```

**New Service Names**:
```
conjugate-master
conjugate-coordination
conjugate-data
```

**DNS Migration**:
```bash
# Update DNS CNAME records
search.example.com → quidditch-coordination.quidditch.svc.cluster.local  # Old
search.example.com → conjugate-coordination.conjugate.svc.cluster.local  # New
```

---

### 3. Docker Images

**Old Images**:
```
docker.io/yourorg/quidditch-master:1.0.0
docker.io/yourorg/quidditch-coordination:1.0.0
docker.io/yourorg/quidditch-data:1.0.0
```

**New Images**:
```
docker.io/yourorg/conjugate-master:1.0.0
docker.io/yourorg/conjugate-coordination:1.0.0
docker.io/yourorg/conjugate-data:1.0.0
```

**Transition Period**: Both image names will be published for 2 months.

---

### 4. Go Module Path

**Old Import Path**:
```go
import (
    "github.com/yourorg/quidditch/pkg/coordination"
    "github.com/yourorg/quidditch/pkg/data"
    "github.com/yourorg/quidditch/pkg/ppl"
)
```

**New Import Path**:
```go
import (
    "github.com/yourorg/conjugate/pkg/coordination"
    "github.com/yourorg/conjugate/pkg/data"
    "github.com/yourorg/conjugate/pkg/ppl"
)
```

**Migration Steps**:
```bash
# Find and replace in your codebase
find . -name "*.go" -type f -exec sed -i 's|github.com/yourorg/quidditch|github.com/yourorg/conjugate|g' {} +

# Update go.mod
go mod edit -replace github.com/yourorg/quidditch=github.com/yourorg/conjugate
go mod tidy
```

---

### 5. Python Pipelines

**Old Python Package**:
```python
from quidditch.pipeline import Processor
from quidditch.client import Client

client = Client("http://localhost:9200")
```

**New Python Package**:
```python
from conjugate.pipeline import Processor
from conjugate.client import Client

client = Client("http://localhost:9200")
```

**Migration Steps**:
```bash
# Update Python dependencies
pip uninstall quidditch
pip install conjugate

# Update imports in your code
find . -name "*.py" -type f -exec sed -i 's/from quidditch/from conjugate/g' {} +
find . -name "*.py" -type f -exec sed -i 's/import quidditch/import conjugate/g' {} +
```

---

### 6. Configuration Files

**Old Config** (`quidditch.yaml`):
```yaml
cluster:
  name: quidditch-cluster
nodes:
  master:
    - quidditch-master-0
    - quidditch-master-1
    - quidditch-master-2
```

**New Config** (`conjugate.yaml`):
```yaml
cluster:
  name: conjugate-cluster
nodes:
  master:
    - conjugate-master-0
    - conjugate-master-1
    - conjugate-master-2
```

---

### 7. CLI Tools

**Old Commands**:
```bash
quidditch cluster status
quidditch index create my-index
quidditch pipeline deploy my-pipeline
```

**New Commands**:
```bash
conjugate cluster status
conjugate index create my-index
conjugate pipeline deploy my-pipeline
```

**Alias (Temporary)**:
```bash
# Add to ~/.bashrc during transition
alias quidditch='conjugate'
```

---

## Common Migration Scenarios

### Scenario 1: Local Development

**Current State**: Running Quidditch locally for development

**Migration Steps**:
1. Pull latest changes: `git pull origin main`
2. Rebuild containers: `docker-compose down && docker-compose up -d`
3. Update environment variables
4. No data migration needed (dev data)

**Estimated Time**: 15 minutes

---

### Scenario 2: Kubernetes Dev Environment

**Current State**: Quidditch deployed in dev Kubernetes namespace

**Migration Steps**:
1. Deploy new CONJUGATE cluster in parallel namespace
2. Validate functionality
3. Switch traffic (update DNS/ingress)
4. Delete old cluster
5. No data migration needed (can start fresh)

**Estimated Time**: 1-2 hours

---

### Scenario 3: Kubernetes Production Environment

**Current State**: Quidditch in production with live traffic

**Migration Steps**:
⚠️ **DO NOT MIGRATE YET** - Wait for Week 4 official procedure

Planned approach:
1. Deploy CONJUGATE cluster in parallel
2. Set up data replication
3. Run both clusters simultaneously
4. Gradual traffic migration (10% → 50% → 100%)
5. Verify data consistency
6. Decommission old cluster

**Estimated Time**: 4-8 hours (with rollback capability)

---

## Rollback Procedures

If you encounter issues during migration:

### Kubernetes Rollback

```bash
# Rollback to previous CRD version
kubectl rollout undo deployment/conjugate-coordination -n conjugate

# Restore old cluster
kubectl apply -f old-cluster-backup.yaml

# Switch DNS back
kubectl patch ingress conjugate-ingress -p '{"spec":{"rules":[{"host":"search.example.com","http":{"paths":[{"path":"/","backend":{"serviceName":"quidditch-coordination","servicePort":9200}}]}}]}}'
```

### Application Rollback

```bash
# Revert Go module path
go mod edit -replace github.com/yourorg/conjugate=github.com/yourorg/quidditch

# Revert Git remote
git remote set-url origin https://github.com/yourorg/quidditch.git
```

---

## Testing Your Migration

### Pre-Migration Checklist

- [ ] Backup all production data
- [ ] Document current configuration
- [ ] Test migration in dev environment
- [ ] Verify all integrations still work
- [ ] Plan rollback procedure
- [ ] Schedule maintenance window

### Post-Migration Verification

```bash
# 1. Check cluster health
conjugate cluster status

# 2. Verify data integrity
curl http://localhost:9200/_cat/indices?v

# 3. Run smoke tests
curl -X POST "http://localhost:9200/test-index/_search" \
  -H 'Content-Type: application/json' \
  -d '{"query": {"match_all": {}}}'

# 4. Check logs for errors
kubectl logs -n conjugate -l app=conjugate-coordination --tail=100

# 5. Monitor metrics
# Check Grafana/Prometheus for anomalies
```

---

## Support During Migration

### Getting Help

- **GitHub Issues**: [Technical problems](https://github.com/yourorg/conjugate/issues)
- **GitHub Discussions**: [Migration questions](https://github.com/yourorg/conjugate/discussions/migration)
- **Slack**: #conjugate-migration channel (coming soon)
- **Email**: migration-support@conjugate.io (coming soon)

### Reporting Issues

If you encounter migration problems:

1. Check this guide for known issues
2. Search existing GitHub issues
3. Create new issue with template:

```markdown
**Migration Issue**: [Brief description]

**Environment**:
- Old version: Quidditch v1.0.0
- New version: CONJUGATE v1.0.0
- Deployment: Kubernetes / Docker / Local
- OS: Linux / Mac / Windows

**Steps to Reproduce**:
1. ...
2. ...

**Error Messages**:
```
[paste error logs]
```

**Expected Behavior**: ...
**Actual Behavior**: ...
```

---

## FAQ

### Q: Do I have to migrate immediately?

**A**: No, but we recommend migrating during the transition period (Jan 29 - Apr 15, 2026) while both names are supported.

### Q: Will my existing data be affected?

**A**: No, data format and storage remain unchanged. Only names and API paths are changing.

### Q: Can I continue using "Quidditch" internally?

**A**: Technically yes, but we recommend adopting the new name to avoid confusion and benefit from future updates.

### Q: What happens after the deprecation deadline (Apr 15, 2026)?

**A**: The old "Quidditch" names will no longer be recognized. You must migrate before this date.

### Q: Will there be breaking API changes?

**A**: No breaking changes to core APIs. Only naming changes to CRDs, packages, and command-line tools.

### Q: How long will the transition period last?

**A**: 10 weeks (Jan 29 - Apr 15, 2026) with full backward compatibility.

### Q: What if I find a bug specific to the migration?

**A**: Report it immediately via GitHub issues with the `migration-bug` label. We'll prioritize fixes.

### Q: Will you provide automated migration scripts?

**A**: Yes, we're developing automated migration tools (coming in Week 3).

### Q: Can I run both Quidditch and CONJUGATE in parallel?

**A**: Yes, during the transition period you can run both clusters side-by-side for testing.

### Q: What about third-party integrations (Grafana, Kibana, etc.)?

**A**: Most integrations use the OpenSearch API, which remains unchanged. Only update dashboard names and labels.

---

## Migration Tools (Coming Soon)

### Automated Migration Script

```bash
# Coming in Week 3
./scripts/migrate-to-conjugate.sh \
  --namespace quidditch \
  --cluster-name quidditch-prod \
  --dry-run

# Actual migration
./scripts/migrate-to-conjugate.sh \
  --namespace quidditch \
  --cluster-name quidditch-prod \
  --backup-data \
  --parallel-deployment
```

### Configuration Converter

```bash
# Coming in Week 3
./tools/convert-config.sh quidditch.yaml > conjugate.yaml
```

### Import Path Updater

```bash
# Coming in Week 3
./tools/update-imports.sh ./my-project
```

---

## Timeline Summary

| Week | Date Range | Activities | Status |
|------|------------|------------|--------|
| **Week 1** | Jan 29 - Feb 4 | Core documentation update, NAMING.md, MIGRATION.md | ✅ In Progress |
| **Week 2** | Feb 5 - Feb 11 | Code updates, CRD changes, package renames | ⏳ Planned |
| **Week 3** | Feb 12 - Feb 18 | Infrastructure updates, migration tools, GitHub rename | ⏳ Planned |
| **Week 4** | Feb 19 - Feb 25 | Public announcement, production migration guide | ⏳ Planned |
| **Weeks 5-10** | Feb 26 - Apr 15 | Transition period, backward compatibility maintained | ⏳ Planned |
| **Apr 15, 2026** | Deadline | End of backward compatibility, old names deprecated | ⏳ Planned |

---

## Resources

- **[NAMING.md](NAMING.md)** - Complete naming rationale
- **[README.md](README.md)** - Updated project overview
- **[CONJUGATE_ARCHITECTURE.md](CONJUGATE_ARCHITECTURE.md)** - Technical architecture (coming soon)
- **[GitHub Discussions - Migration](https://github.com/yourorg/conjugate/discussions/migration)** - Community support

---

## Staying Updated

### Subscribe to Updates

- **GitHub Watch**: Click "Watch" → "Custom" → "Discussions" on the repository
- **GitHub Releases**: Subscribe to release notifications
- **Mailing List**: migration-updates@conjugate.io (coming soon)
- **Twitter**: @conjugate_search (coming soon)

### Release Notes

All migration-related changes will be clearly marked in release notes:
- 🔄 **MIGRATION**: Changes that require action
- 📝 **DEPRECATED**: Old APIs still work but will be removed
- ⚠️ **BREAKING**: Immediate action required (rare)

---

## Conclusion

We understand that migrations can be disruptive. We're committed to making this transition as smooth as possible:

✅ **10-week transition period** with full backward compatibility
✅ **Automated migration tools** (coming Week 3)
✅ **Parallel deployment support** for zero-downtime migration
✅ **Comprehensive documentation** and support
✅ **Clear rollback procedures** if issues arise

The migration to **CONJUGATE** positions us for long-term success by eliminating legal risk and establishing a professional, technically-meaningful brand identity.

**Questions?** Open a discussion on [GitHub Discussions - Migration](https://github.com/yourorg/conjugate/discussions/migration)

---

**Document Version**: 1.0
**Last Updated**: January 29, 2026
**Status**: Active Migration Period
**Target Completion**: February 15, 2026

---

*Made with ❤️ by the CONJUGATE team*
