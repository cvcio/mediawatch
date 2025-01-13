package docker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
)

type Container struct {
	ID   string
	Port string
}

// Destroy stops and removes the container.
func (c *Container) Destroy() error {
	if err := exec.Command("docker", "stop", c.ID).Run(); err != nil {
		return err
	}

	if err := exec.Command("docker", "rm", c.ID, "-v").Run(); err != nil {
		return err
	}

	return nil
}

// StartMongo runs a mongo container to execute commands.
func StartMongo(log *log.Logger) (*Container, error) {
	cmd := exec.Command("docker", "run", "-P", "-d", "mongo:3-jessie")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("starting container: %v", err)
	}

	id := out.String()[:12]

	cmd = exec.Command("docker", "inspect", id)
	out.Reset()
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("inspect container: %v", err)
	}

	var doc []struct {
		NetworkSettings struct {
			Ports struct {
				TCP27017 []struct {
					HostPort string `json:"HostPort"`
				} `json:"27017/tcp"`
			} `json:"Ports"`
		} `json:"NetworkSettings"`
	}
	if err := json.Unmarshal(out.Bytes(), &doc); err != nil {
		return nil, fmt.Errorf("decoding json: %v", err)
	}

	c := Container{
		ID:   id,
		Port: doc[0].NetworkSettings.Ports.TCP27017[0].HostPort,
	}

	return &c, nil
}

// StartPostgres runs a Postgres container.
// func StartPostgres(user, password string) (*Container, error) {
// 	pgUser := fmt.Sprintf("POSTGRES_USER=%s", user)
// 	pgPassword := fmt.Sprintf("POSTGRES_PASSWORD=%s", password)
// 	cmd := exec.Command("docker", "run", "-P",
// 		"-e", pgPassword,
// 		"-e", pgUser,
// 		"-d", "postgres:11")
// 	var out bytes.Buffer
// 	cmd.Stdout = &out
// 	if err := cmd.Run(); err != nil {
// 		return nil, fmt.Errorf("starting container: %v", err)
// 	}
//
// 	id := out.String()[:12]
//
// 	cmd = exec.Command("docker", "inspect", id)
// 	out.Reset()
// 	cmd.Stdout = &out
// 	if err := cmd.Run(); err != nil {
// 		return nil, fmt.Errorf("inspect container: %v", err)
// 	}
//
// 	var doc []struct {
// 		NetworkSettings struct {
// 			Ports struct {
// 				TCP5432 []struct {
// 					HostPort string `json:"HostPort"`
// 				} `json:"5432/tcp"`
// 			} `json:"Ports"`
// 		} `json:"NetworkSettings"`
// 	}
// 	if err := json.Unmarshal(out.Bytes(), &doc); err != nil {
// 		return nil, fmt.Errorf("decoding json: %v", err)
// 	}
//
// 	c := Container{
// 		ID:   id,
// 		Port: doc[0].NetworkSettings.Ports.TCP5432[0].HostPort,
// 	}
//
// 	return &c, nil
// }
