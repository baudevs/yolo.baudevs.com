# 🚀 Contributing to YOLO CLI

First off, thanks for taking the time to contribute! You're awesome, and we want your experience to be great.

## 🎯 Development Process

### Local Development

1. **Fork & Clone**
   ```bash
   git clone https://github.com/yourusername/yolo.baudevs.com.git
   cd yolo.baudevs.com
   ```

2. **Install Dependencies**
   ```bash
   go mod download
   go mod tidy
   ```

3. **Run Tests**
   ```bash
   go test ./...
   ```

4. **Local Installation**
   ```bash
   # Build and install locally
   bash install.sh
   ```

### 🏗️ Project Structure

```
yolo.baudevs.com/
├── cmd/                # Command-line interface
│   └── yolo/          # Main YOLO CLI application
├── internal/          # Internal packages
│   ├── ai/           # AI integration
│   ├── commands/     # CLI commands
│   ├── core/         # Core functionality
│   ├── git/          # Git operations
│   ├── messages/     # Message templates
│   ├── server/       # Web server
│   ├── shortcuts/    # System shortcuts
│   └── web/          # Web interface
├── scripts/          # Build and utility scripts
└── web/             # Web assets
```

### 📝 Coding Guidelines

1. **Code Style**
   - Follow standard Go conventions
   - Use meaningful variable names
   - Write descriptive comments
   - Keep functions focused and small

2. **Documentation**
   - Document all exported functions
   - Include examples in doc strings
   - Keep README.md up to date
   - Add comments for complex logic

3. **Testing**
   - Write unit tests for new features
   - Ensure existing tests pass
   - Add integration tests when needed
   - Test all personality modes

### 🎭 Personality System

YOLO has three personality modes that affect message output:

1. **Clean & Nerdy (1)**
   - Professional and technical
   - Safe for work
   - Light nerdy humor

2. **Mildly Eccentric (2)**
   - More casual tone
   - Occasional sass
   - Pop culture references

3. **Unhinged & Funny (3)**
   - Full chaos mode
   - Heavy on humor
   - Maximum personality

When adding new messages:
- Implement all three personality variants
- Keep the core meaning consistent
- Test with different modes

### 🔄 Git Workflow

1. **Branching**
   ```bash
   # Create a feature branch
   git checkout -b feature/amazing-feature
   
   # Or for bugfixes
   git checkout -b fix/nasty-bug
   ```

2. **Committing**
   ```bash
   # Use YOLO's smart commit
   yolo commit
   
   # Or manually following conventional commits
   git commit -m "feat: add amazing feature"
   ```

3. **Pull Requests**
   - Create PR against `main` branch
   - Fill out the PR template
   - Link related issues
   - Add tests and documentation

### 🐛 Bug Reports

When filing a bug report, please include:

1. YOLO CLI version (`yolo version`)
2. Operating system and version
3. Steps to reproduce
4. Expected vs actual behavior
5. Relevant logs or error messages

### 🚀 Feature Requests

For feature requests:

1. Check existing issues first
2. Describe the feature clearly
3. Explain the use case
4. Suggest implementation if possible
5. Consider all personality modes

## 🤝 Code Review Process

1. **Before Review**
   - Run tests locally
   - Update documentation
   - Check code style
   - Self-review your changes

2. **During Review**
   - Be respectful and constructive
   - Explain your suggestions
   - Link to relevant docs
   - Consider edge cases

3. **After Review**
   - Address all comments
   - Update tests if needed
   - Squash commits if requested
   - Thank your reviewers

## 📜 License

By contributing, you agree that your contributions will be licensed under the MIT License.

## 🎉 Recognition

All contributors are listed in our README.md. We value every contribution, no matter how small!

---

Remember: YOLO isn't just a CLI tool, it's a vibe. Keep it fun, keep it smart, and most importantly, keep it YOLO! 🚀
