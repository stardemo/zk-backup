package zookeeper

import (
	"errors"
	"fmt"
	"github.com/go-zookeeper/zk"
	"log"
	"net"
	"path"
	"strings"
)

func resolveIPv4Addr(addr string) (string, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return "", err
	}
	ipAddrs, err := net.LookupIP(host)
	for _, ipAddr := range ipAddrs {
		ipv4 := ipAddr.To4()
		if ipv4 != nil {
			return net.JoinHostPort(ipv4.String(), port), nil
		}
	}
	return "", fmt.Errorf("no IPv4addr for name %v", host)
}

func resolveZkAddr(zkAddr string) ([]string, error) {
	if zkAddr == "" {
		return nil, nil
	}
	parts := strings.Split(zkAddr, ",")
	resolved := make([]string, 0, len(parts))
	for _, part := range parts {
		// The zookeeper client cannot handle IPv6 addresses before version 3.4.x.
		if r, err := resolveIPv4Addr(part); err != nil {
			log.Printf("cannot resolve %v, will not use it: %v\n", part, err)
		} else {
			resolved = append(resolved, r)
		}
	}
	if len(resolved) == 0 {
		return nil, fmt.Errorf("no valid address found in %v", zkAddr)
	}
	return resolved, nil
}

func DialZk(zkAddr string) (*zk.Conn, <-chan zk.Event, error) {
	resolvedZkAddr, err := resolveZkAddr(zkAddr)
	if err != nil {
		return nil, nil, err
	}

	zconn, session, err := zk.Connect(resolvedZkAddr, 5e9)
	if err == nil {
		// Wait for connection, possibly forever
		event := <-session
		if event.State != zk.StateConnected && event.State != zk.StateConnecting {
			err = fmt.Errorf("zk connect failed: %v", event.State)
		}
		if err == nil {
			return zconn, session, nil
		} else {
			zconn.Close()
		}
	}
	return nil, nil, err
}

func CreateRecursive(conn *zk.Conn, zkPath string, value []byte, flags int32, aclSets []zk.ACL) (pathCreated string, err error) {
	exists, _, err := conn.Exists(zkPath)
	if err != nil {
		return "", err
	}
	if exists {
		conn.Delete(zkPath, -1)
	}
	pathCreated, err = conn.Create(zkPath, value, flags, aclSets)
	if errors.Is(err, zk.ErrNoNode) {
		dirAclv := make([]zk.ACL, len(aclSets))
		for i, acl := range aclSets {
			dirAclv[i] = acl
			dirAclv[i].Perms = PermDirectory
		}
		_, err = CreateRecursive(conn, path.Dir(zkPath), []byte(""), flags, dirAclv)
		if err != nil && !errors.Is(err, zk.ErrNodeExists) {
			return "", err
		}
		pathCreated, err = conn.Create(zkPath, value, flags, aclSets)
	}
	return pathCreated, nil
}
func IsPathExcluded(paths []string, path string) bool {
	for _, p := range paths {
		if p == path {
			return true
		}
	}
	return false
}

func Walk(root string, srcConn *zk.Conn, targetConn *zk.Conn, excludePaths []string) {
	log.Println(srcConn.State())
	children, _, err := srcConn.Children(root)
	if err != nil {
		log.Fatalf("error, when get children of %s, %s\n", root, err)
	}
	for _, node := range children {
		fullPath := path.Join(root, node)
		if IsPathExcluded(excludePaths, fullPath) {
			return
		}
		data, stat, _ := srcConn.Get(fullPath)
		if stat.EphemeralOwner == 0 {
			// ignore ephemeral node
			if _, err := CreateRecursive(targetConn, fullPath, data, 0, zk.WorldACL(zk.PermAll)); err != nil {
				log.Fatalf("error, when create node in target zk, %v\n", err)
			}
			log.Printf("%s backup success", fullPath)
		}
		Walk(fullPath, srcConn, targetConn, excludePaths)
	}
}
