package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

const apktoolVersion = "2.6.0"

func downloadFile(url, outputPath string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP response error: %s", response.Status)
	}

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func decompileAPK(apkFilePath, outputDir string) error {
	cmd := exec.Command("java", "-jar", "apktool.jar", "d", apkFilePath, "-o", outputDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("APK decompilation failed: %v", err)
	}

	fmt.Println("APK decompilation completed.")
	return nil
}

func main() {
	// Replace with the path to your APK file
	apkFilePath := "com.goibibo.apk"

	// Replace with the desired output directory
	outputDir := "Day3_Task"

	// Download Apktool JAR from an alternative mirror
	err := downloadFile(fmt.Sprintf("https://github.com/iBotPeaches/Apktool/releases/download/v%s/apktool_%s.jar", apktoolVersion, apktoolVersion), "apktool.jar")
	if err != nil {
		fmt.Println("Error downloading Apktool JAR:", err)
		os.Exit(1)
	}

	// Perform the APK decompilation
	err = decompileAPK(apkFilePath, outputDir)
	if err != nil {
		fmt.Println("Error during APK decompilation:", err)
		os.Exit(1)
	}
}
