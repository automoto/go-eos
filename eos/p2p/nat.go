package p2p

import (
	"context"
	"fmt"

	"github.com/mydev/go-eos/eos/internal/callback"
	"github.com/mydev/go-eos/eos/internal/cbinding"
	"github.com/mydev/go-eos/eos/types"
)

// NATType describes the NAT class observed for the local network. Players
// behind a Strict NAT can typically only connect to Open or Moderate peers,
// so games often store this as a lobby attribute (see PRD LOB-14) for
// matchmaking compatibility filtering.
type NATType int

const (
	// NATUnknown indicates the NAT type has not been queried yet.
	NATUnknown NATType = 0
	// NATOpen indicates an open NAT with no connection restrictions.
	NATOpen NATType = 1
	// NATModerate indicates a moderate NAT that can connect to most peers.
	NATModerate NATType = 2
	// NATStrict indicates a strict NAT that can typically only connect to open or moderate peers.
	NATStrict NATType = 3
)

// String returns a human-readable label for the NAT type.
func (n NATType) String() string {
	switch n {
	case NATOpen:
		return "Open"
	case NATModerate:
		return "Moderate"
	case NATStrict:
		return "Strict"
	default:
		return "Unknown"
	}
}

// QueryNATType performs an active probe to determine the local network's
// NAT type. The result is also cached and retrievable via GetNATType.
//
// This is the only async operation in the P2P interface — everything else
// is synchronous or notification-based.
func (p *P2P) QueryNATType(ctx context.Context) (NATType, error) {
	oneshot := callback.NewOneShot()
	if err := p.worker.Submit(func() {
		cbinding.EOS_P2P_QueryNATType(p.handle, oneshot.HandleValue())
	}); err != nil {
		oneshot.Delete()
		return NATUnknown, err
	}

	result, err := oneshot.Wait(ctx)
	if err != nil {
		return NATUnknown, err
	}

	info := result.Data.(*cbinding.EOS_P2P_OnQueryNATTypeCompleteInfo)
	if info.ResultCode != cbinding.EOS_EResult_Success {
		return NATUnknown, types.NewResult(int(info.ResultCode))
	}
	return NATType(info.NATType), nil
}

// GetNATType returns the most recently queried NAT type. Returns NATUnknown
// and an *eos/types.Result wrapping EOS_NotFound if QueryNATType has never
// been called successfully.
func (p *P2P) GetNATType() (NATType, error) {
	var natType cbinding.EOS_ENATType
	var result cbinding.EOS_EResult
	if err := p.worker.Submit(func() {
		natType, result = cbinding.EOS_P2P_GetNATType(p.handle)
	}); err != nil {
		return NATUnknown, fmt.Errorf("p2p get nat: %w", err)
	}
	if result != cbinding.EOS_EResult_Success {
		return NATUnknown, types.NewResult(int(result))
	}
	return NATType(natType), nil
}

// RelayControl mirrors EOS_ERelayControl. Compatibility matrix:
//
//	NoRelays    + ForceRelays  → incompatible
//	NoRelays    + (No|Allow)   → compatible
//	AllowRelays + everything   → compatible (default)
//	ForceRelays + AllowRelays  → compatible
type RelayControl int

const (
	// NoRelays disables relay servers; only direct connections are used.
	NoRelays RelayControl = 0
	// AllowRelays allows relay servers as a fallback (default).
	AllowRelays RelayControl = 1
	// ForceRelays forces all traffic through relay servers.
	ForceRelays RelayControl = 2
)

// SetRelayControl configures whether peer connections may use Epic relay
// servers. Applies to new and renegotiating connections, not existing ones.
func (p *P2P) SetRelayControl(rc RelayControl) error {
	var result cbinding.EOS_EResult
	if err := p.worker.Submit(func() {
		result = cbinding.EOS_P2P_SetRelayControl(p.handle, &cbinding.EOS_P2P_SetRelayControlOptions{
			RelayControl: cbinding.EOS_ERelayControl(rc),
		})
	}); err != nil {
		return fmt.Errorf("p2p set relay: %w", err)
	}
	if result != cbinding.EOS_EResult_Success {
		return types.NewResult(int(result))
	}
	return nil
}

// GetRelayControl returns the current relay control setting.
func (p *P2P) GetRelayControl() (RelayControl, error) {
	var rc cbinding.EOS_ERelayControl
	var result cbinding.EOS_EResult
	if err := p.worker.Submit(func() {
		rc, result = cbinding.EOS_P2P_GetRelayControl(p.handle)
	}); err != nil {
		return AllowRelays, fmt.Errorf("p2p get relay: %w", err)
	}
	if result != cbinding.EOS_EResult_Success {
		return AllowRelays, types.NewResult(int(result))
	}
	return RelayControl(rc), nil
}

// PortRange describes the port window the SDK uses for outgoing P2P traffic.
type PortRange struct {
	Port                    uint16
	MaxAdditionalPortsToTry uint16
}

// SetPortRange configures the port range. Defaults are 7777 + 99 additional.
// Set Port=0 (and MaxAdditionalPortsToTry=0) to let the OS choose.
func (p *P2P) SetPortRange(r PortRange) error {
	var result cbinding.EOS_EResult
	if err := p.worker.Submit(func() {
		result = cbinding.EOS_P2P_SetPortRange(p.handle, &cbinding.EOS_P2P_SetPortRangeOptions{
			Port:                    r.Port,
			MaxAdditionalPortsToTry: r.MaxAdditionalPortsToTry,
		})
	}); err != nil {
		return fmt.Errorf("p2p set port range: %w", err)
	}
	if result != cbinding.EOS_EResult_Success {
		return types.NewResult(int(result))
	}
	return nil
}

// GetPortRange returns the current port range.
func (p *P2P) GetPortRange() (PortRange, error) {
	var port, maxAdditional uint16
	var result cbinding.EOS_EResult
	if err := p.worker.Submit(func() {
		port, maxAdditional, result = cbinding.EOS_P2P_GetPortRange(p.handle)
	}); err != nil {
		return PortRange{}, fmt.Errorf("p2p get port range: %w", err)
	}
	if result != cbinding.EOS_EResult_Success {
		return PortRange{}, types.NewResult(int(result))
	}
	return PortRange{Port: port, MaxAdditionalPortsToTry: maxAdditional}, nil
}

// PacketQueueInfo is the diagnostic snapshot of the SDK's packet queues.
type PacketQueueInfo struct {
	IncomingMaxSizeBytes     uint64
	IncomingCurrentSizeBytes uint64
	IncomingCurrentPackets   uint64
	OutgoingMaxSizeBytes     uint64
	OutgoingCurrentSizeBytes uint64
	OutgoingCurrentPackets   uint64
}

// GetPacketQueueInfo returns current packet queue diagnostics. Use to detect
// queue pressure (e.g., before deciding to drop low-priority packets).
func (p *P2P) GetPacketQueueInfo() (PacketQueueInfo, error) {
	var info *cbinding.EOS_P2P_PacketQueueInfo
	var result cbinding.EOS_EResult
	if err := p.worker.Submit(func() {
		info, result = cbinding.EOS_P2P_GetPacketQueueInfo(p.handle)
	}); err != nil {
		return PacketQueueInfo{}, fmt.Errorf("p2p packet queue info: %w", err)
	}
	if result != cbinding.EOS_EResult_Success {
		return PacketQueueInfo{}, types.NewResult(int(result))
	}
	return PacketQueueInfo{
		IncomingMaxSizeBytes:     info.IncomingPacketQueueMaxSizeBytes,
		IncomingCurrentSizeBytes: info.IncomingPacketQueueCurrentSizeBytes,
		IncomingCurrentPackets:   info.IncomingPacketQueueCurrentPacketCount,
		OutgoingMaxSizeBytes:     info.OutgoingPacketQueueMaxSizeBytes,
		OutgoingCurrentSizeBytes: info.OutgoingPacketQueueCurrentSizeBytes,
		OutgoingCurrentPackets:   info.OutgoingPacketQueueCurrentPacketCount,
	}, nil
}

// SetPacketQueueSize sets the maximum byte sizes for the incoming and
// outgoing packet queues. 0 means unlimited.
func (p *P2P) SetPacketQueueSize(incomingMax, outgoingMax uint64) error {
	var result cbinding.EOS_EResult
	if err := p.worker.Submit(func() {
		result = cbinding.EOS_P2P_SetPacketQueueSize(p.handle, &cbinding.EOS_P2P_SetPacketQueueSizeOptions{
			IncomingPacketQueueMaxSizeBytes: incomingMax,
			OutgoingPacketQueueMaxSizeBytes: outgoingMax,
		})
	}); err != nil {
		return fmt.Errorf("p2p set queue size: %w", err)
	}
	if result != cbinding.EOS_EResult_Success {
		return types.NewResult(int(result))
	}
	return nil
}

// ClearPacketQueue drops all queued (incoming + outgoing) packets between
// localUserId and remoteUserId on the given socket.
func (p *P2P) ClearPacketQueue(localUserId, remoteUserId types.ProductUserId, socket SocketId) error {
	if err := socket.Validate(); err != nil {
		return err
	}
	cSocket := &cbinding.EOS_P2P_SocketId{Name: socket.Name}
	var result cbinding.EOS_EResult
	if err := p.worker.Submit(func() {
		result = cbinding.EOS_P2P_ClearPacketQueue(p.handle, &cbinding.EOS_P2P_ClearPacketQueueOptions{
			LocalUserId:  cbinding.EOS_ProductUserId_FromString(string(localUserId)),
			RemoteUserId: cbinding.EOS_ProductUserId_FromString(string(remoteUserId)),
			SocketId:     cSocket,
		})
	}); err != nil {
		return fmt.Errorf("p2p clear queue: %w", err)
	}
	if result != cbinding.EOS_EResult_Success {
		return types.NewResult(int(result))
	}
	return nil
}
