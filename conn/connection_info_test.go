package conn

//func TestMarshalConnInfo(t *testing.T) {
//
//	km, err := key.NewKeyManager("~/key")
//
//	defer os.RemoveAll("~/key")
//
//	_, pub, err := km.GenerateKey(key.RSA4096)
//
//	connInfo := NewConnInfo(FromPubKey(pub), Address{IP: "127.0.0.1:8888"}, pub)
//
//	b, err := json.Marshal(connInfo)
//
//	fmt.Printf("[%s]", b)
//
//	if err != nil {
//
//	}
//
//	connectedConnInfo := &ConnInfo{}
//	err = json.Unmarshal(b, connectedConnInfo)
//
//	assert.NoError(t, err)
//}

//func TestHostInfo_GetPublicInfo(t *testing.T) {
//
//	//given
//	km, err := key.NewKeyManager("~/key")
//	assert.NoError(t, err)
//
//	defer os.RemoveAll("~/key")
//
//	pri, pub, err := km.GenerateKey(key.RSA4096)
//	assert.NoError(t, err)
//
//	hostInfo := NewHostInfo(FromPubKey(pub), Address{IP: "127.0.0.1:8888"}, pub, pri)
//
//	//when
//	pInfo := hostInfo.GetPublicInfo()
//
//	//then
//	assert.NotNil(t, pInfo)
//	assert.Equal(t, hostInfo)
//	log.Print(pInfo)
//}
