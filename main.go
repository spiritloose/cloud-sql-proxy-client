package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"time"
)

const mysqlBinEnv = "CLOUD_SQL_PROXY_CLIENT_MYSQL"

const mysqlBinDefault = "mysql"

const cloudSQLProxyBin = "cloud_sql_proxy"

const cloudSQLProxyWaitTimeout = 10

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: cloud-sql-proxy-client INSTANCE [MYSQL_ARGS...]")
		return
	}
	instance := os.Args[1]
	tempdir, err := ioutil.TempDir("", "cloud-sql-proxy-client-")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	proxy, err := launchCloudSQLProxy(instance, tempdir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.RemoveAll(tempdir)
		return
	}
	socket := path.Join(tempdir, instance)
	exist := waitSocketFile(socket)
	if !exist {
		fmt.Fprintln(os.Stderr, "socket file not found")
		proxy.Process.Kill()
		proxy.Wait()
		os.RemoveAll(tempdir)
		return
	}

	err = runMySQL(socket)
	exitCode := 0
	if err != nil {
		exitCode = 1
	}

	proxy.Process.Kill()
	proxy.Wait()
	os.RemoveAll(tempdir)

	os.Exit(exitCode)
}

func runMySQL(socket string) error {
	args := []string{"--socket", socket}
	args = append(args, os.Args[2:]...)
	mysqlBin := os.Getenv(mysqlBinEnv)
	if mysqlBin == "" {
		mysqlBin = mysqlBinDefault
	}
	mysql, err := exec.LookPath(mysqlBin)
	if err != nil {
		return err
	}
	cmd := exec.Command(mysql, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func launchCloudSQLProxy(instance string, dir string) (*exec.Cmd, error) {
	cloudSQL, err := exec.LookPath(cloudSQLProxyBin)
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(cloudSQL, "-dir", dir, "-instances", instance)
	out, err := os.Create(path.Join(dir, "cloud_sql_proxy.log"))
	if err != nil {
		return nil, err
	}
	defer out.Close()
	cmd.Stdout = out
	cmd.Stderr = out
	err = cmd.Start()
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

func waitSocketFile(socket string) bool {
	found := false
	for i := 0; i < cloudSQLProxyWaitTimeout; i++ {
		_, err := os.Stat(socket)
		if err == nil {
			found = true
			break
		}
		time.Sleep(1 * time.Second)
	}
	return found
}
