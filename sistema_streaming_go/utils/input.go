package utils

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var lector = bufio.NewReader(os.Stdin)

func LeerTexto(mensaje string) string {
	fmt.Print(mensaje)
	texto, _ := lector.ReadString('\n')
	return strings.TrimSpace(texto)
}

func LeerEntero(mensaje string) (int, error) {
	texto := LeerTexto(mensaje)
	numero, err := strconv.Atoi(texto)
	if err != nil {
		return 0, fmt.Errorf("debe ingresar un número entero: %w", err)
	}
	return numero, nil
}
