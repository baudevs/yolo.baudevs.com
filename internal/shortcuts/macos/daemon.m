#import <Cocoa/Cocoa.h>
#import <Carbon/Carbon.h>

// Callback function type
typedef void (*ShortcutCallback)(const char* identifier);

// External Go callback
extern void shortcutTriggeredCallback(const char* identifier);

// Global variables
static CFMachPortRef eventTap;
static CFRunLoopSourceRef runLoopSource;
static ShortcutCallback shortcutCallback;
static NSMutableDictionary* registeredHotKeys;

// Forward declarations
static void createEventTap(void);

// Initialize on load
__attribute__((constructor)) static void initialize(void) {
    registeredHotKeys = [[NSMutableDictionary alloc] init];
}

// Set callback function
void setShortcutCallback(ShortcutCallback callback) {
    shortcutCallback = callback;
}

// Event tap callback
static CGEventRef eventTapCallback(CGEventTapProxy proxy, CGEventType type, CGEventRef event, void *refcon) {
    if (type != kCGEventKeyDown) {
        return event;
    }

    // Get key code and modifiers
    CGKeyCode keyCode = (CGKeyCode)CGEventGetIntegerValueField(event, kCGKeyboardEventKeycode);
    CGEventFlags flags = CGEventGetFlags(event);

    // Convert to Carbon modifiers
    int modifiers = 0;
    if (flags & kCGEventFlagMaskCommand) modifiers |= cmdKey;
    if (flags & kCGEventFlagMaskShift) modifiers |= shiftKey;
    if (flags & kCGEventFlagMaskAlternate) modifiers |= optionKey;
    if (flags & kCGEventFlagMaskControl) modifiers |= controlKey;

    // Check if this matches any registered shortcuts
    NSString* identifier = nil;
    for (NSString* key in registeredHotKeys) {
        NSDictionary* hotKey = registeredHotKeys[key];
        if ([hotKey[@"keyCode"] intValue] == keyCode && 
            [hotKey[@"modifiers"] intValue] == modifiers) {
            identifier = key;
            break;
        }
    }

    // If found, trigger callback
    if (identifier && shortcutCallback) {
        shortcutCallback([identifier UTF8String]);
        return NULL; // Consume the event
    }

    return event;
}

// Helper function to create the event tap
static void createEventTap(void) {
    // Create event tap
    eventTap = CGEventTapCreate(
        kCGSessionEventTap,
        kCGHeadInsertEventTap,
        kCGEventTapOptionDefault,
        CGEventMaskBit(kCGEventKeyDown),
        eventTapCallback,
        NULL
    );

    if (!eventTap) {
        NSLog(@"Failed to create event tap - please check accessibility permissions");
        return;
    }

    // Create run loop source
    runLoopSource = CFMachPortCreateRunLoopSource(kCFAllocatorDefault, eventTap, 0);
    CFRunLoopAddSource(CFRunLoopGetCurrent(), runLoopSource, kCFRunLoopCommonModes);
    CGEventTapEnable(eventTap, true);

    // Start run loop in a background thread
    dispatch_async(dispatch_get_global_queue(DISPATCH_QUEUE_PRIORITY_DEFAULT, 0), ^{
        CFRunLoopRun();
    });
}

// Start event tap
void startEventTap(void) {
    // Check accessibility permissions first
    if (!AXIsProcessTrusted()) {
        NSLog(@"⚠️ Accessibility permissions required");
        NSLog(@"Please follow these steps:");
        NSLog(@"1. Open System Settings > Privacy & Security > Accessibility");
        NSLog(@"2. Click the + button");
        NSLog(@"3. Navigate to and select the YOLO application");
        NSLog(@"4. Ensure the checkbox next to YOLO is checked");
        
        // Request permissions with prompt
        NSDictionary *options = @{
            (__bridge NSString *)kAXTrustedCheckOptionPrompt: @YES
        };
        AXIsProcessTrustedWithOptions((__bridge CFDictionaryRef)options);
        
        // Wait a bit to see if permissions were granted
        dispatch_after(dispatch_time(DISPATCH_TIME_NOW, 1 * NSEC_PER_SEC), dispatch_get_main_queue(), ^{
            if (AXIsProcessTrusted()) {
                // Try creating the event tap again
                createEventTap();
            }
        });
        return;
    }
    
    createEventTap();
}

// Stop event tap
void stopEventTap(void) {
    if (runLoopSource) {
        CFRunLoopRemoveSource(CFRunLoopGetCurrent(), runLoopSource, kCFRunLoopCommonModes);
        CFRelease(runLoopSource);
        runLoopSource = NULL;
    }

    if (eventTap) {
        CFRelease(eventTap);
        eventTap = NULL;
    }
}

// Register hot key
void registerHotKey(const char* identifier, int keyCode, int modifiers) {
    NSString* nsIdentifier = [NSString stringWithUTF8String:identifier];
    NSDictionary* hotKey = @{
        @"keyCode": @(keyCode),
        @"modifiers": @(modifiers)
    };
    registeredHotKeys[nsIdentifier] = hotKey;
}

// Unregister hot key
void unregisterHotKey(const char* identifier) {
    NSString* nsIdentifier = [NSString stringWithUTF8String:identifier];
    [registeredHotKeys removeObjectForKey:nsIdentifier];
} 