#!/bin/bash

# Create YOLO directory structure
mkdir -p yolo/{epics,features,tasks}

# Create initial files
cat > yolo/README.md << 'EOF'
# Project YOLO Documentation

This project follows the YOLO methodology for documentation and project management.
See YOLO.md in the root directory for full methodology details.

## Quick Start
1. New features go in `features/`
2. Implementation tasks go in `tasks/`
3. Strategic initiatives go in `epics/`
4. Never delete, only mark as deprecated
5. Always link tasks to epics
6. Keep history in CHANGELOG.md and HISTORY.yml
EOF

# Create initial HISTORY.yml
cat > HISTORY.yml << 'EOF'
version: 0.1.0
date: $(date +%Y-%m-%d)
changes:
  - type: feature
    description: Initial project setup with YOLO methodology
    impact: Established foundation for project documentation
    files:
      - YOLO.md
      - HISTORY.yml
      - CHANGELOG.md
      - yolo/*
EOF

# Create initial CHANGELOG.md
cat > CHANGELOG.md << 'EOF'
# Changelog

## [0.1.0] - $(date +%Y-%m-%d)

### Added
- Initial YOLO methodology implementation
- Basic project structure
- Documentation framework
EOF

# Create LLM instruction file
cat > LLM_INSTRUCTIONS.md << 'EOF'
# Instructions for LLM Developers

This project follows the YOLO methodology (see YOLO.md). Key points:

1. Never delete content, only mark as deprecated/changed
2. All tasks must link to epics via tags
3. Maintain full history in HISTORY.yml and CHANGELOG.md
4. Follow semantic versioning
5. Keep relationships between documents
6. Document all decisions and their context

Directory Structure:
```
yolo/
├── epics/     # Strategic initiatives
├── features/  # Feature specs
└── tasks/     # Implementation tasks
```

When assisting:
1. Read HISTORY.yml for context
2. Check task dependencies
3. Preserve all historical information
4. Add appropriate status tags
5. Update version control files
EOF

# Copy YOLO.md from root if it exists, otherwise download it
if [ -f YOLO.md ]; then
    cp YOLO.md yolo/
else
    echo "Please copy YOLO.md to the project root"
fi

# Make script executable
chmod +x init-yolo.sh

echo "YOLO methodology initialized successfully!"
echo "Next steps:"
echo "1. Review YOLO.md for methodology details"
echo "2. Check LLM_INSTRUCTIONS.md for AI developer guidelines"
echo "3. Start creating epics in yolo/epics/"
echo "4. Link tasks to epics in yolo/tasks/" 