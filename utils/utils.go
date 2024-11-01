package utils

import "os/exec"



func IsGitInstalled() (bool, error) {
    _, err := exec.LookPath("git")
    return err != nil, err
}
