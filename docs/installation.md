# Installing YOLO: A Beginner's Guide

This guide will walk you through installing YOLO on your computer, even if you've never used a command line before.

## What is the Command Line?

Before we start, let's understand what the command line is:
- It's a text-based way to interact with your computer
- On Mac, it's called "Terminal"
- On Windows, it's called "Command Prompt" or "PowerShell"

## Step-by-Step Installation

### Step 1: Opening the Terminal

**On Mac:**
1. Press `Command (⌘)` + `Space` to open Spotlight
2. Type "Terminal"
3. Click on the Terminal app or press Enter

**On Windows:**
1. Press `Windows Key` + `R`
2. Type "powershell"
3. Press Enter

### Step 2: Installing YOLO

We've made installation super simple with our install script. Just copy and paste the following command:

```bash
curl -fsSL https://yolo.baudevs.com/install.sh | sh
```

Here's what this does:
1. Downloads our installation script
2. Runs it automatically
3. Sets up YOLO on your computer

If you see any messages asking for your password, this is normal. Type your computer's password and press Enter (you won't see the characters as you type - this is for security).

### Step 3: Verifying the Installation

To make sure YOLO was installed correctly:

1. Type this command and press Enter:
   ```bash
   yolo --version
   ```

2. You should see something like:
   ```
   YOLO version 1.2.3
   ```

### Step 4: First-Time Setup

Let's configure YOLO for first use:

1. Run the setup wizard:
   ```bash
   yolo setup
   ```

2. The wizard will ask you some questions:
   - Your name (for commit messages)
   - Your email
   - If you have an OpenAI API key (optional)
   - If you have a YOLO license key (optional)

Just type your answers and press Enter after each one.

### Common Questions

#### Q: What if I see "command not found"?
This means your computer needs to restart the Terminal. Close it and open a new Terminal window.

#### Q: What if I get a permission error?
Try adding `sudo` before the command:
```bash
sudo curl -fsSL https://yolo.baudevs.com/install.sh | sh
```

#### Q: How do I know it worked?
Try a simple YOLO command:
```bash
yolo hello
```
You should see a friendly greeting!

### Next Steps

Now that YOLO is installed, try these beginner-friendly commands:

1. Get help:
   ```bash
   yolo help
   ```

2. Create your first task:
   ```bash
   yolo task create "My first task"
   ```

3. See what YOLO can do:
   ```bash
   yolo suggest
   ```

## Troubleshooting

### If the Installation Fails

1. Try running this first:
   ```bash
   xcode-select --install
   ```
   This installs required developer tools on Mac.

2. Make sure you're connected to the internet

3. If you still have problems, run:
   ```bash
   yolo doctor
   ```
   This will check your system and tell you what's wrong.

### Getting Help

If you're stuck:

1. Visit our help site: [help.baudevs.com](https://help.baudevs.com)
2. Join our community: [community.baudevs.com](https://community.baudevs.com)
3. Email support: support@baudevs.com

## System Requirements

- **Mac:** macOS 10.15 or newer
- **Windows:** Windows 10 or newer
- **Linux:** Ubuntu 20.04 or newer
- Internet connection
- 500MB free disk space

## Uninstalling

If you need to remove YOLO:

```bash
yolo uninstall
```

This will cleanly remove YOLO from your system.

## Updating YOLO

To get the latest version:

```bash
yolo update
```

YOLO will automatically check for updates and install them.

Remember: YOLO is designed to be user-friendly! If you're ever unsure about something, just type `yolo help` or ask our community. We're here to help you succeed!
