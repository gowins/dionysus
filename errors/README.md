# pkg说明

### errors

此包主要作用是聚合errors,将多个不同的error聚合成一个error。如果是同样的error则只会被记录一次。
使用方式如下所示：
```golang
func closeClients(rdclis []*redis.Client) error {
	var errlist []error
	for _, rdcli := range rdclis {
		if rdcli != nil {
			err := rdcli.Close()
			if err != nil {
				errlist = append(errlist, err)
			}
		}
	}
	return utilerrors.NewAggregate(errlist)
}
```