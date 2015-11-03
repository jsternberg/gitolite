package gitolite

import "golang.org/x/crypto/ssh"

type channelWriter struct {
	channel ssh.Channel
}

func (w *channelWriter) Write(data []byte) (int, error) {
	return w.channel.Write(data)
}

func (w *channelWriter) Close() error {
	return w.channel.CloseWrite()
}

type channelReader struct {
	channel ssh.Channel
}

func (r *channelReader) Read(data []byte) (int, error) {
	return r.channel.Read(data)
}

func (r *channelReader) Close() error {
	return r.channel.Close()
}
