package utils

import (
	"fmt"
	"net"
	"os"
	"reflect"
	"regexp"

	"github.com/BurntSushi/toml"
	"github.com/kiracore/sekin/src/shidai/internal/types"
	"github.com/tyler-smith/go-bip39"
	"go.uber.org/zap"
)

// ValidateIP checks if the given string is a valid IPv4 or IPv6 address.
// It returns true if the IP is valid, otherwise returns false.
func ValidateIP(ip string) bool {
	isValid := net.ParseIP(ip) != nil
	zap.L().Debug("Validating IP", zap.String("IP", ip), zap.Bool("IsValid", isValid))
	return isValid
}

// ValidatePort checks if the given port number is within the valid range of 1 to 65535.
// It returns true if the port is valid, otherwise returns false.
func ValidatePort(port int) bool {
	isValid := port > 0 && port <= 65535
	zap.L().Debug("Validating port", zap.Int("Port", port), zap.Bool("IsValid", isValid))
	return isValid
}

// ValidateMnemonic checks if the given mnemonic is valid according to the BIP-0039 standard.
// It returns true if the mnemonic is valid, otherwise returns false.
func ValidateMnemonic(mnemonic string) bool {
	isValid := bip39.IsMnemonicValid(mnemonic)
	zap.L().Debug("Validating mnemonic", zap.String("Mnemonic", mnemonic), zap.Bool("IsValid", isValid))
	return isValid
}

// IsPublicIP checks if the given IP address is a public IP address.
func IsPublicIP(ip net.IP) bool {
	privateIPBlocks := []*regexp.Regexp{
		regexp.MustCompile(`^10\..*`),
		regexp.MustCompile(`^172\.(1[6-9]|2[0-9]|3[0-1])\..*`),
		regexp.MustCompile(`^192\.168\..*`),
	}
	ipStr := ip.String()
	for _, block := range privateIPBlocks {
		if block.MatchString(ipStr) {
			zap.L().Debug("IP identified as private", zap.String("IP", ipStr))
			return false
		}
	}
	return true
}

// GetPublicIP retrieves the public IP address of the system.
// Returns an error if more than one public IP address is found.
func GetPublicIP() (string, error) {
	var publicIPs []string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		zap.L().Error("Failed to get interface addresses", zap.Error(err))
		return "", fmt.Errorf("failed to get interface addresses: %w", err)
	}

	for _, addr := range addrs {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}

		if ip != nil && !ip.IsLoopback() && ip.To4() != nil && IsPublicIP(ip) {
			publicIPs = append(publicIPs, ip.String())
		}
	}

	if len(publicIPs) == 0 {
		zap.L().Error("No public IP addresses found")
		return "", fmt.Errorf("no public IP addresses found")
	}
	if len(publicIPs) > 1 {
		zap.L().Error("Multiple public IP addresses found")
		return "", fmt.Errorf("multiple public IP addresses found")
	}

	zap.L().Info("Public IP address retrieved", zap.String("IP", publicIPs[0]))
	return publicIPs[0], nil
}

func FileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		zap.L().Debug("File does not exist", zap.String("filePath", filePath))
		return false
	}
	isFile := !info.IsDir()
	zap.L().Debug("File existence checked", zap.String("filePath", filePath), zap.Bool("Exists", isFile))
	return isFile
}

// DeleteFile removes a file specified by the file path.
func DeleteFile(filePath string) error {
	// Log the attempt to delete the file
	zap.L().Info("Attempting to delete file", zap.String("filePath", filePath))

	err := os.Remove(filePath)
	if err != nil {
		// Log the error if the file could not be deleted
		zap.L().Error("Failed to delete file", zap.String("filePath", filePath), zap.Error(err))
		return fmt.Errorf("failed to delete file %s: %w", filePath, err)
	}

	// Log successful deletion
	zap.L().Info("File deleted successfully", zap.String("filePath", filePath))
	return nil
}

func CreateDir(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

// LoadConfig loads config.toml to Config structure
func LoadConfig(filePath string, config types.Config) error {
	zap.L().Info("Attempting to load configuration", zap.String("filePath", filePath))
	if _, err := toml.DecodeFile(filePath, config); err != nil {
		zap.L().Error("Failed to load config from file", zap.String("filePath", filePath), zap.Error(err))
		return fmt.Errorf("failed to load config: %w", err)
	}
	zap.L().Info("Configuration loaded successfully", zap.String("filePath", filePath))
	return nil
}

// LoadAppConfig loads app.toml to ConfigApp structure
func LoadAppConfig(filePath string, config types.AppConfig) error {
	zap.L().Info("Attempting to load app configuration", zap.String("filePath", filePath))
	if _, err := toml.DecodeFile(filePath, config); err != nil {
		zap.L().Error("Failed to load app config from file", zap.String("filePath", filePath), zap.Error(err))
		return fmt.Errorf("failed to load app config: %w", err)
	}
	zap.L().Info("App configuration loaded successfully", zap.String("filePath", filePath))
	return nil
}

// SaveConfig saves config.toml to given path
func SaveConfig(filePath string, config types.Config) error {
	zap.L().Info("Attempting to save configuration", zap.String("filePath", filePath))
	file, err := os.Create(filePath)
	if err != nil {
		zap.L().Error("Failed to create config file", zap.String("filePath", filePath), zap.Error(err))
		return fmt.Errorf("failed to create config file: %w", err)
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)
	if err := encoder.Encode(config); err != nil {
		zap.L().Error("Failed to encode configuration to file", zap.String("filePath", filePath), zap.Error(err))
		return fmt.Errorf("failed to encode config: %w", err)
	}
	zap.L().Info("Configuration saved successfully", zap.String("filePath", filePath))
	return nil
}

// SaveAppConfig saves app.toml to given path
func SaveAppConfig(filePath string, config types.AppConfig) error {
	zap.L().Info("Attempting to save app configuration", zap.String("filePath", filePath))
	file, err := os.Create(filePath)
	if err != nil {
		zap.L().Error("Failed to create app config file", zap.String("filePath", filePath), zap.Error(err))
		return fmt.Errorf("failed to create app config file: %w", err)
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)
	if err := encoder.Encode(config); err != nil {
		zap.L().Error("Failed to encode app configuration to file", zap.String("filePath", filePath), zap.Error(err))
		return fmt.Errorf("failed to encode app config: %w", err)
	}
	zap.L().Info("App configuration saved successfully", zap.String("filePath", filePath))
	return nil
}

// SetField sets a value on a given struct field and returns a description of the change.
func SetField(obj interface{}, fieldName string, newValue interface{}) (string, error) {
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		zap.L().Error("Invalid object type", zap.Any("Type", v.Type()))
		return "", fmt.Errorf("expected a pointer to a struct")
	}

	v = v.Elem() // Dereference the pointer to get the struct

	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		zap.L().Error("No such field", zap.String("Field", fieldName))
		return "", fmt.Errorf("no such field: %s in obj", fieldName)
	}
	if !field.CanSet() {
		zap.L().Error("Cannot set field", zap.String("Field", fieldName))
		return "", fmt.Errorf("cannot set field %s", fieldName)
	}

	fieldType := field.Type()
	newVal := reflect.ValueOf(newValue)
	if newVal.Type() != fieldType {
		zap.L().Error("Type mismatch",
			zap.String("Expected", fieldType.String()),
			zap.String("Got", newVal.Type().String()))
		return "", fmt.Errorf("provided value type %s doesn't match obj field type %s", newVal.Type(), fieldType)
	}

	oldValue := fmt.Sprintf("%v", field.Interface())
	field.Set(newVal)
	changeDescription := fmt.Sprintf("Changed %s from %s to %v", fieldName, oldValue, newValue)

	zap.L().Info("Field updated", zap.String("Change", changeDescription))

	return changeDescription, nil
}
