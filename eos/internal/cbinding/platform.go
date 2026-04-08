//go:build !eosstub

package cbinding

/*
#include "platform_wrapper.h"
#include "eos_sdk.h"
#include <stdlib.h>
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
