package service

import (
	"encoding/json"
	"fmt"
)

// Small no-op references so files that import encoding/json or fmt don't fail
var _ = json.Marshal
var _ = fmt.Sprintf
