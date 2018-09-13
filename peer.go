package bifrost

type PeerInfo struct {
	IP       string
	Pubkey   []byte
	CurveOpt int
}
