// ======================================================================================
// cgo compilation (for desktop platforms and local tests)
// ======================================================================================

#include <stdio.h>
#include <stddef.h>
#include <stdbool.h>
#include "_cgo_export.h"

typedef void (*callback)(const char *jsonEvent);
callback gCallback = 0;

bool StatusServiceSignalEvent(const char *jsonEvent) {
	if (gCallback) {
		gCallback(jsonEvent);
	}

	return true;
}

void SetEventCallback(void *cb) {
	gCallback = (callback)cb;
}
