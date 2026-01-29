# Phase 3: Infrastructure Updates Guide

**Date**: January 29, 2026
**Status**: Ready for Execution
**Estimated Duration**: 3-4 hours

---

## Overview

Phase 3 involves updating external infrastructure that cannot be changed programmatically. This guide provides step-by-step instructions for all required manual and automated tasks.

---

## Quick Status

### ✅ Already Complete (Phase 2)

- [x] CI/CD workflow files updated
- [x] Docker build configurations updated
- [x] Deployment scripts updated
- [x] All code references updated

### 📋 Manual Steps Required (This Phase)

- [ ] **GitHub Repository Rename** (15 min)
- [ ] **Docker Hub Configuration** (30 min)
- [ ] **Domain Registration** (1 hour)
- [ ] **DNS Configuration** (30 min)
- [ ] **Update External References** (1 hour)

---

## Section 1: GitHub Repository Rename

### Prerequisites

- Admin access to the GitHub repository
- Backup of all important data
- Communication plan for users

### Steps

#### 1.1 Pre-Rename Checklist

```bash
# 1. Verify current state
git status
git log --oneline -n 5

# 2. Create backup branch
git branch backup-pre-rename

# 3. Tag current state
git tag -a phase2-complete -m "Phase 2: Code Updates Complete"
git push origin phase2-complete

# 4. Ensure all changes are committed
git commit -am "Phase 3: Ready for repository rename"
git push origin main
```

#### 1.2 Rename Repository on GitHub

**URL**: https://github.com/yourorg/quidditch/settings

1. Navigate to repository settings
2. Scroll to "Repository name" section
3. Change name from `quidditch` to `conjugate`
4. Click "Rename"
5. GitHub will automatically create a redirect from the old name

**Result**:
- Old URL: `https://github.com/yourorg/quidditch`
- New URL: `https://github.com/yourorg/conjugate`
- Redirect: Automatic for 1 year

#### 1.3 Update Local Clones

**For all developers and CI/CD systems**:

```bash
# Option 1: Update remote URL
cd /path/to/quidditch
git remote set-url origin https://github.com/yourorg/conjugate.git
git remote -v  # Verify

# Option 2: Fresh clone (recommended for clean start)
cd /path/to/projects
git clone https://github.com/yourorg/conjugate.git
cd conjugate
```

#### 1.4 Update Repository Settings

After rename, update these settings:

**Description**:
```
Cloud-native Observability + Natural-language Joint Understanding Granular search Analytics Tunable Engine - A distributed search platform with OpenSearch compatibility
```

**Topics/Tags**:
```
search-engine, distributed-systems, cloud-native, observability,
nlp, kubernetes, opensearch, diagon, conjugate, analytics
```

**About Section**:
- ✅ Website: `https://conjugate.io` (when ready)
- ✅ Issues enabled
- ✅ Wiki disabled (use docs/ directory)
- ✅ Discussions enabled

#### 1.5 Update Branch Protection Rules

Verify branch protection rules still apply after rename:

```
Settings → Branches → Branch protection rules
```

- `main` branch: Require PR, require status checks
- `develop` branch: Same rules as main

#### 1.6 Verify GitHub Actions

After rename, trigger workflows:

```bash
# Push a small change to trigger CI
git commit --allow-empty -m "Trigger CI after rename"
git push origin main
```

Check: https://github.com/yourorg/conjugate/actions

---

## Section 2: Docker Image Publishing

### Prerequisites

- Docker Hub account or GitHub Container Registry access
- Authentication configured

### Steps

#### 2.1 GitHub Container Registry (GHCR) ✅ Automatic

The Docker workflow (`.github/workflows/docker.yml`) automatically publishes to GHCR:

**Images**:
- `ghcr.io/yourorg/conjugate/master:latest`
- `ghcr.io/yourorg/conjugate/coordination:latest`
- `ghcr.io/yourorg/conjugate/data:latest`

**No action required** - images will be published automatically after repository rename.

#### 2.2 Docker Hub (Optional)

If you want to publish to Docker Hub as well:

##### Step 1: Create Docker Hub Organization

1. Go to https://hub.docker.com
2. Create organization: `conjugate`
3. Create repositories:
   - `conjugate/master`
   - `conjugate/coordination`
   - `conjugate/data`

##### Step 2: Add Docker Hub Secrets

```
GitHub Settings → Secrets and variables → Actions → New repository secret
```

Add:
- `DOCKERHUB_USERNAME`: Your Docker Hub username
- `DOCKERHUB_TOKEN`: Docker Hub access token

##### Step 3: Update Docker Workflow

Add Docker Hub login to `.github/workflows/docker.yml`:

```yaml
- name: Log in to Docker Hub
  uses: docker/login-action@v3
  with:
    username: ${{ secrets.DOCKERHUB_USERNAME }}
    password: ${{ secrets.DOCKERHUB_TOKEN }}

- name: Extract metadata for Docker Hub
  id: meta-dockerhub
  uses: docker/metadata-action@v5
  with:
    images: conjugate/master  # or coordination, data
    tags: |
      type=ref,event=branch
      type=sha,prefix={{branch}}-
      type=raw,value=latest,enable={{is_default_branch}}

# Add to build-and-push tags:
tags: |
  ${{ steps.meta.outputs.tags }}
  ${{ steps.meta-dockerhub.outputs.tags }}
```

##### Step 4: Trigger Build

```bash
git commit -am "Add Docker Hub publishing"
git push origin main
```

#### 2.3 Verify Image Publication

```bash
# Pull from GHCR
docker pull ghcr.io/yourorg/conjugate/master:latest

# Pull from Docker Hub (if configured)
docker pull conjugate/master:latest

# Run quick test
docker run --rm ghcr.io/yourorg/conjugate/master:latest --version
```

#### 2.4 Update Documentation

Update image references in:
- `README.md` - Usage examples
- `MIGRATION.md` - Docker image migration
- `deployments/docker-compose/docker-compose.yml` - Image references
- `deployments/kubernetes/*.yaml` - Image references

---

## Section 3: Domain Registration

### Prerequisites

- Domain registrar account (Namecheap, GoDaddy, Route 53, etc.)
- Credit card for payment
- Email for domain verification

### Steps

#### 3.1 Register Primary Domain

**Recommended**: `conjugate.io`

1. Go to your preferred registrar
2. Search for `conjugate.io`
3. Purchase for 1-3 years
4. Enable auto-renewal
5. Add privacy protection (if available)

**Alternative domains to register**:
- `conjugate.dev` - Developer documentation
- `conjugate.com` - Marketing/commercial (optional)
- `conjugate.cloud` - Cloud platform (optional)

**Estimated Cost**:
- `.io` domain: $30-40/year
- `.dev` domain: $12-15/year
- `.com` domain: $10-15/year

#### 3.2 Configure Name Servers

**Option A: Use Registrar's DNS**
- Use default name servers from registrar
- Configure DNS records in next section

**Option B: Use Cloudflare (Recommended)**
- Create Cloudflare account (free tier)
- Add site: `conjugate.io`
- Update name servers at registrar
- Benefits: Free SSL, CDN, DDoS protection

**Option C: Use AWS Route 53**
- Create hosted zone for `conjugate.io`
- Update name servers at registrar
- Benefits: Integration with AWS services

#### 3.3 Verify Domain Ownership

```bash
# Check DNS propagation
dig conjugate.io NS
dig conjugate.dev NS

# Check with multiple resolvers
nslookup conjugate.io 8.8.8.8
nslookup conjugate.io 1.1.1.1
```

Wait 24-48 hours for full propagation.

---

## Section 4: DNS Configuration

### Prerequisites

- Domain registered and verified
- DNS management access
- SSL certificate plan (Let's Encrypt recommended)

### Steps

#### 4.1 Configure DNS Records

**For `conjugate.io`**:

```
# A records (when hosting is ready)
@               A       <your-server-ip>
www             A       <your-server-ip>
docs            A       <docs-server-ip>
api             A       <api-server-ip>

# CNAME for GitHub Pages (if using)
docs            CNAME   yourorg.github.io.

# MX records (email - optional)
@               MX 10   mx1.forwardemail.net.
@               MX 20   mx2.forwardemail.net.

# TXT records
@               TXT     "v=spf1 include:_spf.forwardemail.net ~all"
```

**For `conjugate.dev`**:

```
# Point to documentation
@               CNAME   yourorg.github.io.
www             CNAME   yourorg.github.io.
```

#### 4.2 SSL/TLS Configuration

**Option A: Let's Encrypt (Recommended)**

```bash
# Using certbot
sudo certbot certonly --standalone \
  -d conjugate.io \
  -d www.conjugate.io \
  -d docs.conjugate.io \
  -d api.conjugate.io
```

**Option B: Cloudflare (Free)**

- Enable "Full (strict)" SSL mode
- Automatic HTTPS Rewrites: On
- Always Use HTTPS: On
- Automatic certificate provisioning

#### 4.3 Verify DNS Configuration

```bash
# Test DNS resolution
dig conjugate.io A
dig www.conjugate.io A
dig docs.conjugate.io CNAME

# Test SSL
curl -I https://conjugate.io
openssl s_client -connect conjugate.io:443

# Check SSL certificate
echo | openssl s_client -connect conjugate.io:443 2>/dev/null | openssl x509 -noout -dates
```

#### 4.4 Update Repository Settings

**GitHub Pages** (if using for docs):

```
Settings → Pages
→ Custom domain: docs.conjugate.io
→ Enforce HTTPS: ✅
```

---

## Section 5: Update External References

### Prerequisites

- Access to all external platforms
- List of all references to update

### Steps

#### 5.1 Documentation Sites

**If using GitHub Pages**:

```bash
# In docs repository
echo "docs.conjugate.io" > CNAME
git add CNAME
git commit -m "Add custom domain"
git push origin gh-pages
```

#### 5.2 Package Registries

**Go Module Proxy**:
- No action needed - Go modules use repository URL

**Docker Hub** (if publishing):
- Update organization description
- Update image descriptions
- Add links to new repository

#### 5.3 Social Media / Forums

Update references on:
- **Twitter/X**: Profile, pinned posts
- **Reddit**: r/golang posts
- **Hacker News**: User profile
- **Stack Overflow**: Tags, posts
- **LinkedIn**: Project posts

#### 5.4 Blog Posts / Articles

If you've published articles:
- Update author bio links
- Add redirect notices
- Update code examples

#### 5.5 Dependencies

Check projects that depend on this:
```bash
# Search GitHub for dependents
# https://github.com/yourorg/conjugate/network/dependents
```

Notify maintainers of module path change.

#### 5.6 Update README Badges

```markdown
<!-- Update badges -->
[![Status](https://img.shields.io/badge/status-beta-yellow)](https://conjugate.io)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/yourorg/conjugate)](https://goreportcard.com/report/github.com/yourorg/conjugate)
[![GoDoc](https://godoc.org/github.com/yourorg/conjugate?status.svg)](https://godoc.org/github.com/yourorg/conjugate)
[![Docker Pulls](https://img.shields.io/docker/pulls/conjugate/master.svg)](https://hub.docker.com/r/conjugate/master/)
```

---

## Section 6: CI/CD Updates

### Prerequisites

- CI/CD pipelines identified
- Access to CI/CD platforms

### Steps

#### 6.1 GitHub Actions ✅ Already Updated

GitHub Actions automatically uses the new repository name via `${{ github.repository }}`.

**Verify workflows**:
```bash
# Check workflow status
https://github.com/yourorg/conjugate/actions
```

#### 6.2 External CI Systems

If using Travis CI, CircleCI, Jenkins, etc.:

**Travis CI**:
1. Go to travis-ci.com
2. Sync repositories
3. Verify webhook points to new URL

**CircleCI**:
1. Go to circleci.com
2. Project Settings → Advanced Settings
3. Update GitHub repository URL

**Jenkins**:
1. Update Git repository URL in job configuration
2. Test webhook integration

#### 6.3 Monitoring Integration

**Codecov**:
```bash
# Update repository in Codecov dashboard
https://codecov.io/gh/yourorg/conjugate/settings
```

**SonarCloud**:
```bash
# Re-import project if needed
https://sonarcloud.io/projects/create
```

---

## Section 7: Communication & Announcement

### Prerequisites

- Draft announcement ready
- Communication channels identified

### Steps

#### 7.1 GitHub Release

Create release announcement:

```markdown
# CONJUGATE v1.0.0 - Official Rebranding

## 🎉 Major Update: Project Rebranded

The project formerly known as "Quidditch" has been officially rebranded to **CONJUGATE**.

**New Name**: CONJUGATE
**Backronym**: Cloud-native Observability + Natural-language Joint Understanding Granular search Analytics Tunable Engine

## Why the Change?

- **Legal**: Avoid trademark conflicts
- **Professional**: Establish technical identity
- **Strategic**: Emphasize four core pillars

## What Changed?

- ✅ Project name: Quidditch → CONJUGATE
- ✅ Repository: yourorg/quidditch → yourorg/conjugate
- ✅ Go module: github.com/yourorg/quidditch → github.com/yourorg/conjugate
- ✅ Docker images: quidditch/* → conjugate/*
- ✅ All documentation updated

## Migration Guide

See [MIGRATION.md](MIGRATION.md) for complete migration instructions.

## Links

- 🏠 Website: https://conjugate.io
- 📖 Documentation: https://docs.conjugate.io
- 💬 Discussions: https://github.com/yourorg/conjugate/discussions
```

#### 7.2 GitHub Discussions

Post announcement:

```markdown
Title: Project Rebranded to CONJUGATE

We're excited to announce the official rebranding of our project to **CONJUGATE**!

[Same content as release notes]

## Questions?

Please ask in this thread or open an issue.
```

#### 7.3 Email Notification

If you have a mailing list:

```
Subject: Important: Quidditch → CONJUGATE Rebranding

Dear Community,

We're writing to inform you of an important update...

[Same content as release notes]

Best regards,
The CONJUGATE Team
```

#### 7.4 Update Documentation

Add banner to old documentation (if still hosted):

```html
<div class="alert alert-warning">
  ⚠️ This project has been rebranded to CONJUGATE.
  Please visit <a href="https://conjugate.io">conjugate.io</a> for updated documentation.
</div>
```

---

## Section 8: Verification Checklist

### Post-Migration Verification

- [ ] **GitHub Repository**
  - [ ] Repository renamed successfully
  - [ ] Redirect working from old URL
  - [ ] CI/CD workflows passing
  - [ ] Issues/PRs accessible

- [ ] **Docker Images**
  - [ ] Images published to GHCR
  - [ ] Images pullable with new names
  - [ ] Tags correct (latest, versioned)

- [ ] **Domains**
  - [ ] conjugate.io registered
  - [ ] DNS configured
  - [ ] SSL certificate active
  - [ ] Website accessible

- [ ] **Documentation**
  - [ ] README updated
  - [ ] All links working
  - [ ] Badges updated
  - [ ] Examples using new names

- [ ] **External References**
  - [ ] Social media updated
  - [ ] Blog posts updated
  - [ ] Package registries updated

- [ ] **Communication**
  - [ ] GitHub release published
  - [ ] Announcement in Discussions
  - [ ] Email sent (if applicable)

---

## Section 9: Rollback Plan

If critical issues arise:

### GitHub Repository Rollback

```bash
# Rename back to quidditch (temporary)
# GitHub Settings → Repository name → quidditch

# Update local clones
git remote set-url origin https://github.com/yourorg/quidditch.git
```

### Docker Images Rollback

```bash
# Old images remain available during transition
docker pull ghcr.io/yourorg/quidditch/master:latest
```

### DNS Rollback

```bash
# Remove new DNS records
# Point domain to parking page
# Refund domain (if within grace period)
```

---

## Section 10: Timeline

| Task | Duration | Dependencies | Status |
|------|----------|--------------|--------|
| **GitHub Rename** | 15 min | None | ⏳ Pending |
| **Update Local Clones** | 5 min | GitHub rename | ⏳ Pending |
| **Docker Hub Setup** | 30 min | GitHub rename | ⏳ Pending |
| **Domain Registration** | 1 hour | None (can parallelize) | ⏳ Pending |
| **DNS Configuration** | 30 min | Domain registration | ⏳ Pending |
| **SSL Setup** | 30 min | DNS configuration | ⏳ Pending |
| **External References** | 1 hour | GitHub rename | ⏳ Pending |
| **Announcement** | 30 min | All above complete | ⏳ Pending |
| **Verification** | 30 min | Announcement sent | ⏳ Pending |
| **TOTAL** | **4-5 hours** | | |

---

## Section 11: Success Criteria

Phase 3 is complete when:

- [x] GitHub repository renamed to `conjugate`
- [x] Docker images published as `conjugate/*`
- [x] Domain `conjugate.io` registered and configured
- [x] DNS resolving correctly
- [x] SSL certificate active
- [x] All CI/CD workflows passing
- [x] Documentation updated
- [x] Public announcement made
- [x] No broken links or references

---

## Section 12: Support & Resources

### Getting Help

- **GitHub Issues**: Technical problems
- **GitHub Discussions**: General questions
- **Email**: (to be set up)

### Useful Links

- **NAMING.md**: Complete naming rationale
- **MIGRATION.md**: User migration guide
- **PHASE2_REBRANDING_COMPLETE.md**: Code updates report
- **REBRANDING_SUMMARY.md**: Overall status

---

## Conclusion

Phase 3 involves primarily manual steps for external infrastructure. Follow this guide sequentially, verify each step, and use the rollback plan if needed.

**Estimated Time**: 4-5 hours
**Complexity**: Medium (mostly administrative tasks)
**Risk**: Low (can rollback if needed)

**Next Phase**: Phase 4 (Communication & Community Migration)

---

**Document Version**: 1.0
**Last Updated**: January 29, 2026
**Status**: Ready for Execution

---

*Made with ❤️ by the CONJUGATE team*
