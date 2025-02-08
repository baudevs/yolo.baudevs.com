# YOLO Methodology: A Fun & Straightforward Approach to AI Project Evolution

Welcome to the “YOLO” methodology! In this guide, we’ll explore how YOLO helps human developers (even complete beginners) create AI-enhanced projects without fear of losing track, deleting crucial history, or getting bogged down by complicated processes. We’ll explain all the core ideas, folder structures, and prompts you might need. By the end, you’ll see how easy and enjoyable it can be to maintain a complete log of your code’s evolution.

---

## 1. What Is YOLO?

YOLO stands for **“You Observe, Log, and Oversee.”** It’s a simple but powerful approach to AI-driven software development. In essence:  

- You (the developer) observe and record each coding step you take, including changes, additions, or removals.  
- You keep a continuous log of these changes in friendly, readable files (like `HISTORY.yml` and `CHANGELOG.md`).  
- You oversee the entire project’s progress by marking each change as “deprecated,” “added,” “removed,” etc., instead of fully deleting it.

This ensures a permanent record of how your project evolves—making it easier to learn, retrace steps, undo mistakes, or figure out why something was done a certain way.

---

## 2. Why YOLO?

Most coding methodologies for advanced AI projects can be overwhelming. YOLO simplifies it so that human developers of all experience levels can:  
1. Keep a day-by-day (or step-by-step) history of changes.  
2. Never lose important context owing to deleted lines or code.  
3. Make it fun by eliminating the dryness of typical version control rules.  
4. Use straightforward folder structures and naming conventions that new developers can follow.  
5. Adapt easily to different AI tools and LLMs.

---

## 3. Where Do We Put Our Files?

We encourage you to keep a folder called `yolo/` in your project. This folder will store your epics, tasks, relationships, and any additional files that detail how your project is evolving. Beyond that, you typically want to maintain these files in the root of your project:

- **`README.md`**: A plain-English introduction to your project.  
- **`CHANGELOG.md`**: A chronological summary of your additions, changes, and fixes, following a simple version numbering (SemVer).  
- **`WISHES.md`**: A fun place to capture future ideas or features you “wish” to add someday.  
- **`STRATEGY.md`** (optional): A short outline of your current approach or game plan.  
- **`HISTORY.yml`**: A YAML file that tracks every notable action in a structured, date-tagged format.

A typical project layout might look like this:

```bash
my-cool-project/
├── yolo/
│   ├── epics/
│   ├── features/
│   ├── tasks/
│   └── ...
├── README.md
├── CHANGELOG.md
├── WISHES.md
├── STRATEGY.md
├── HISTORY.yml
└── src/
    └── ...
```

---

## 4. How the Logging Works

### 4.1 No Line Is Truly Removed  
In YOLO, we don’t delete lines from logs. If you decide a piece of code or text no longer belongs, mark it as **“deprecated”** or **“removed”** in the logs. That way, you can always trace back how the code used to be.

### 4.2 Quick Tags and Changes  
When updating your code or text, explain **why** so future-you (or collaborators) can see the rationale. You’ll commonly use tags like:  
- `added`: for something newly introduced  
- `deprecated`: for something you plan to remove or no longer consider best practice  
- `removed`: to mark lines or features you have ended  
- `changed`: for updates or revisions  
- `implemented`: when something is fully functional  
- `not relevant`: for lines that no longer matter

### 4.3 Keep It Easy  
Remember, YOLO is aimed at new developers trying AI coding for the first time. Keep your notes in plain, friendly language. For example:  
```yaml
- type: feature
  description: "Added new AI prompt generation script"
  impact: "Simplifies daily coding tasks"
  status: "implemented"
```

---

## 5. Example of `HISTORY.yml` Logging

Below is a simplified example of how your `HISTORY.yml` might look. Notice that this is for demonstration—your actual file can vary based on what you do each day:

```yaml
version: 1.2.0
date: 2025-01-13
changes:
  - type: feature
    description: "Introduced advanced AI code completion"
    impact: "Speeds up coding time for new developers"
    files:
      - src/ai-completion.js
    status: "implemented"

  - type: docs
    description: "Created user-friendly YOLO manual"
    impact: "Helps newcomers quickly learn YOLO"
    files:
      - README.md
      - yolo/intro.md
    status: "added"

  - type: deprecated
    description: "Deprecated old code snippet for text analysis"
    impact: "We have better AI functions now"
    files:
      - src/legacy-text-analysis.js
    status: "deprecated"
```

Notice how we never “delete” an old snippet from the file; we just mark it as deprecated or removed. That helps keep a complete history.

---

## 6. The Core Prompt for YOLO Methodology

When using an AI assistant (e.g., ChatGPT or similar LLM), after you create or update code, you can instruct the AI to log changes by using a standardized prompt. For example:

“awesome! now document everything according to our yolo methodology in @yolo folder adding the epics, tasks and the relationships between all of them, update the @README.md @WISHES.md @STRATEGY.md and @HISTORY.yml and @CHANGELOG.md as we do in the YOLO methodology. Remember: we don't delete anything or lines, we just mark as deprecated, added, implemented, not relevant, anything that changes, so we can keep a complete history of the evolution of our project. CHANGELOG.md and README.md maintain standard development practices, like SemVer.”

This prompt tells the LLM that it should:  
1. Update the relevant files in the `yolo` folder.  
2. Keep track of new epics, tasks, or relationships.  
3. Append new entries to `HISTORY.yml` (instead of deleting old entries).  
4. Maintain semver rules in `CHANGELOG.md`.  
5. Never remove lines from logs—only mark them accordingly.

---

## 7. Example Use Cases

1. **Adding a New Feature**  
   - You write the code for a new AI-based playlist generator.  
   - You run your code, it works, and you’re ready to log it.  
   - You say:  
     “awesome! now document everything according to our yolo methodology …”  
   - The AI updates your `CHANGELOG.md`, `HISTORY.yml`, and also creates a mention in `@README.md` about the new feature.

2. **Deprecating an Old Function**  
   - You realize an old function is no longer needed.  
   - Instead of removing it from the logs entirely, you mark it as “deprecated” in `HISTORY.yml` and `CHANGELOG.md`.  
   - The function can physically remain in the code if you want, or you can remove the function from the code while keeping a note in the logs (depicting it as “removed”).

---

## 8. Future Enhancements to the Philosophy

Although YOLO is already fun and beginner-friendly, we have plans to improve it further:

1. **Automated LLM Checklists**  
   - Generate a daily “to-do” list from `WISHES.md` and `STRATEGY.md`.  
   - Prompt the developer if items remain incomplete.

2. **Built-In Visualization**  
   - Transform `HISTORY.yml` into a timeline or flowchart automatically.  
   - Show your project’s evolution at a glance.

3. **Metadata Hooks**  
   - Tag lines with optional metadata (like bug tracking ID or important references).  
   - Summaries can be generated to track tasks across epics in a single view.

4. **Beginner Tutorials**  
   - Step-by-step “first commit” tutorials.  
   - LLM-based guided onboarding that references real code in your project.

---

## 9. Concluding Thoughts on YOLO

By adopting the YOLO methodology, you free yourself from the anxiety of losing crucial lines of code or historical context. Instead, you build up a story of how your AI-driven project grows—every small success, every pivot, every moment of creativity stays in the record. YOLO helps you:

- Keep the entire evolution of your project transparent.  
- Make it simple and fun to develop with AI.  
- Provide easy references for new and experienced contributors alike.  
- Offer a stepping stone for novices who have never coded until now, but can do so confidently with AI by their side.

---

## 10. Quick Reference

| File           | Purpose                                                   |
|----------------|-----------------------------------------------------------|
| `yolo/`        | Contains epics, tasks, and relationship docs about your project evolution. |
| `README.md`    | High-level summary of your project (what it is, how to use it).             |
| `CHANGELOG.md` | List of changes by version (semver style). Always appended, never truncated.|
| `HISTORY.yml`  | YAML file with day-to-day or action-by-action records of modifications.     |
| `WISHES.md`    | A simple file for brainstorming or “wish list” items you’d like to see.     |
| `STRATEGY.md`  | (Optional) Outlines your project’s current direction or approach.           |

Remember:  
1. **Never** remove lines from logs—just mark them accordingly if they change status.  
2. Use simple, descriptive language to keep it clear for all developers.  
3. Always keep track of your entire journey, not just the highlights!

---

### Final Word

Welcome aboard the YOLO journey. Embrace the freedom and creativity it provides, and enjoy building AI projects without the usual stress of losing context. Keep it fun, keep it simple, and log everything in a way that tells your project’s story—from the first tiny improvement to your biggest breakthrough.

Happy developing with YOLO!
