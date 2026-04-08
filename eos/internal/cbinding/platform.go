//go:build !eosstub

package cbinding

/*
#cgo CFLAGS: -I${SRCDIR}/../../../static/EOS-SDK-49960398-Release-v1.19.0.3/SDK/Include
#cgo darwin LDFLAGS: -L${SRCDIR}/../../../static/EOS-SDK-49960398-Release-v1.19.0.3/SDK/Bin -lEOSSDK-Mac-Shipping -Wl,-rpath,${SRCDIR}/../../../static/EOS-SDK-49960398-Release-v1.19.0.3/SDK/Bin
#cgo linux LDFLAGS: -L${SRCDIR}/../../../static/EOS-SDK-49960398-Release-v1.19.0.3/SDK/Bin -lEOSSDK-Linux-Shipping -Wl,-rpath,${SRCDIR}/../../../static/EOS-SDK-49960398-Release-v1.19.0.3/SDK/Bin

#include "eos_sdk.h"
#include "eos_init.h"
#include "eos_logging.h"
#include <stdlib.h>
#include <stdint.h>

static EOS_EResult eos_initialize(const char* productName, const char* productVersion) {
    EOS_InitializeOptions opts = {0};
    opts.ApiVersion = EOS_INITIALIZE_API_LATEST;
    opts.ProductName = productName;
    opts.ProductVersion = productVersion;
    return EOS_Initialize(&opts);
}

static uintptr_t eos_platform_create(const char* productId, const char* sandboxId,
    const char* deploymentId, const char* clientId, const char* clientSecret) {
    EOS_Platform_Options opts = {0};
    opts.ApiVersion = EOS_PLATFORM_OPTIONS_API_LATEST;
    opts.ProductId = productId;
    opts.SandboxId = sandboxId;
    opts.DeploymentId = deploymentId;
    opts.ClientCredentials.ClientId = clientId;
    opts.ClientCredentials.ClientSecret = clientSecret;
    opts.bIsServer = EOS_FALSE;
    return (uintptr_t)EOS_Platform_Create(&opts);
}

static void eos_platform_tick(uintptr_t handle) {
    EOS_Platform_Tick((EOS_HPlatform)handle);
}

static void eos_platform_release(uintptr_t handle) {
    EOS_Platform_Release((EOS_HPlatform)handle);
}

static uintptr_t eos_platform_get_auth(uintptr_t handle) {
    return (uintptr_t)EOS_Platform_GetAuthInterface((EOS_HPlatform)handle);
}

static uintptr_t eos_platform_get_connect(uintptr_t handle) {
    return (uintptr_t)EOS_Platform_GetConnectInterface((EOS_HPlatform)handle);
}

static uintptr_t eos_platform_get_lobby(uintptr_t handle) {
    return (uintptr_t)EOS_Platform_GetLobbyInterface((EOS_HPlatform)handle);
}

static uintptr_t eos_platform_get_sessions(uintptr_t handle) {
    return (uintptr_t)EOS_Platform_GetSessionsInterface((EOS_HPlatform)handle);
}

static uintptr_t eos_platform_get_p2p(uintptr_t handle) {
    return (uintptr_t)EOS_Platform_GetP2PInterface((EOS_HPlatform)handle);
}
*/
import "C"
import "unsafe"

type EOS_InitializeOptions struct {
	ProductName    string
	ProductVersion string
}

type EOS_Platform_Options struct {
	ProductId    string
	SandboxId    string
	DeploymentId string
	ClientId     string
	ClientSecret string
}

func EOS_Initialize(opts *EOS_InitializeOptions) EOS_EResult {
	cName := C.CString(opts.ProductName)
	defer C.free(unsafe.Pointer(cName))
	cVersion := C.CString(opts.ProductVersion)
	defer C.free(unsafe.Pointer(cVersion))

	return EOS_EResult(C.eos_initialize(cName, cVersion))
}

func EOS_Platform_Create(opts *EOS_Platform_Options) EOS_HPlatform {
	cProductId := C.CString(opts.ProductId)
	defer C.free(unsafe.Pointer(cProductId))
	cSandboxId := C.CString(opts.SandboxId)
	defer C.free(unsafe.Pointer(cSandboxId))
	cDeploymentId := C.CString(opts.DeploymentId)
	defer C.free(unsafe.Pointer(cDeploymentId))
	cClientId := C.CString(opts.ClientId)
	defer C.free(unsafe.Pointer(cClientId))
	cClientSecret := C.CString(opts.ClientSecret)
	defer C.free(unsafe.Pointer(cClientSecret))

	return EOS_HPlatform(C.eos_platform_create(cProductId, cSandboxId, cDeploymentId, cClientId, cClientSecret))
}

func EOS_Platform_Tick(handle EOS_HPlatform) {
	C.eos_platform_tick(C.uintptr_t(handle))
}

func EOS_Platform_Release(handle EOS_HPlatform) {
	C.eos_platform_release(C.uintptr_t(handle))
}

func EOS_Shutdown() EOS_EResult {
	return EOS_EResult(C.EOS_Shutdown())
}

func EOS_Platform_GetAuthInterface(handle EOS_HPlatform) EOS_HAuth {
	return EOS_HAuth(C.eos_platform_get_auth(C.uintptr_t(handle)))
}

func EOS_Platform_GetConnectInterface(handle EOS_HPlatform) EOS_HConnect {
	return EOS_HConnect(C.eos_platform_get_connect(C.uintptr_t(handle)))
}

func EOS_Platform_GetLobbyInterface(handle EOS_HPlatform) EOS_HLobby {
	return EOS_HLobby(C.eos_platform_get_lobby(C.uintptr_t(handle)))
}

func EOS_Platform_GetSessionsInterface(handle EOS_HPlatform) EOS_HSessions {
	return EOS_HSessions(C.eos_platform_get_sessions(C.uintptr_t(handle)))
}

func EOS_Platform_GetP2PInterface(handle EOS_HPlatform) EOS_HP2P {
	return EOS_HP2P(C.eos_platform_get_p2p(C.uintptr_t(handle)))
}
