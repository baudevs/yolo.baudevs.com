# YOLO Documentation

## Project Structure

### Epics
- [E001] Interactive Project Initialization
  - Status: In Progress
  - Features: F001-F004, F011
  - Description: Interactive project setup and configuration

- [E002] Project Visualization
  - Status: Implemented
  - Features: F005-F008
  - Description: 3D visualization of project relationships

### Features
- [F001] Interactive TUI Setup ✓
- [F002] Project Detection and Reconfiguration ✓
- [F003] Git Integration ✓
- [F004] AI Provider Integration ✓
- [F005] 3D Graph Visualization ✓
- [F006] Interactive Controls ✓
- [F007] Real-time Updates ✓
- [F008] Embedded Web Server ✓
- [F011] System-wide Keyboard Shortcuts ⚡
  - Status: In Progress
  - Web Interface: Implemented ✓
  - Configuration: Implemented ✓
  - macOS Daemon: Pending

### Tasks
- [T001-T005] Project Initialization Tasks ✓
- [T010-T017] Graph Visualization Tasks ✓
- [T011.1-T011.5] Shortcuts Configuration Tasks
  - T011.1: Web Interface ✓
  - T011.2: macOS Daemon ⚡
  - T011.3: Configuration Command ✓
  - T011.4: Persistence Layer ✓
  - T011.5: Management Features ✓

### Relationships
- [R001] Project Initialization to Graph Visualization
  - Type: Dependency
  - Status: Active
- [R002] Shortcuts to AI Integration
  - Type: Enhancement
  - Status: Active
  - Description: Enables quick access to AI features via shortcuts

## Documentation Guidelines
1. Never delete information, mark as:
   - Deprecated
   - Superseded
   - Not relevant
   - Implemented
   - Updated
   - Improved
   - Fixed

2. Always add context:
   - Creation dates
   - Update history
   - Implementation details
   - Technical notes
   - Related items

3. Maintain relationships:
   - Between epics
   - Features to epics
   - Tasks to features
   - Cross-component dependencies

4. Version control:
   - Follow SemVer
   - Update CHANGELOG.md
   - Track in HISTORY.yml
   - Document migrations

## Legend
✓ - Implemented
⚡ - In Progress