package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"sync"

	"github.com/fatih/color"
)

var foundlist []string
var mutex sync.Mutex

func main() {
	banner := `
   ██████╗  ██████╗ ███████╗██╗   ██╗███████╗███████╗
  ██╔════╝ ██╔═══██╗██╔════╝██║   ██║╚══███╔╝╚══███╔╝
  ██║  ███╗██║   ██║█████╗  ██║   ██║  ███╔╝   ███╔╝ 
  ██║   ██║██║   ██║██╔══╝  ██║   ██║ ███╔╝   ███╔╝  
  ╚██████╔╝╚██████╔╝██║     ╚██████╔╝███████╗███████╗
   ╚═════╝  ╚═════╝ ╚═╝      ╚═════╝ ╚══════╝╚══════╝
                  by @RaiVivashu`

	color.Red(banner)
	var url string
	var extension string
	var wordlist string
	var threads int
	var outputfile string
	threads = 5

	//users flags
	color.Blue("\nEnter the URL:\n")
	fmt.Scan(&url)
	color.Blue("\nEnter the extension you want to check:\n")
	fmt.Scan(&extension)
	color.Blue("\nEnter the wordlist path:\n")
	fmt.Scan(&wordlist)
	color.Blue("\nEnter the number of threads to use\n")
	fmt.Scan(&threads)
	color.Blue("\nEnter the Output file\n")
	fmt.Scan(&outputfile)
	fmt.Println("\n\n")
	//open wordlists file
	file, err := os.Open(wordlist)
	if err != nil {
		color.Red("Wrong File Provided!!!")
	}

	//export all the wordlists
	scanner := bufio.NewScanner(file)
	var words []string
	for scanner.Scan() {
		word := scanner.Text()
		words = append(words, word)
	}
	defer file.Close()
	//export all the extensions
	reg, error := regexp.Compile(`\b\w+\b`)
	if error != nil {
		color.Red("\nPlease Enter a valid extensions like php,txt,html etc....\n")
	}
	extensions := reg.FindAllString(extension, -1)

	//fuzzing starts
	var wg sync.WaitGroup

	t1 := make(map[string]bool)
	t2 := make(map[string]bool)
	for i := 1; i <= threads; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			//mutex.Lock()
			//defer mutex.Unlock()
			t1 = fuzz_entensions(words, url, extensions, t1)
		}()
		go func() {
			defer wg.Done()
			//mutex.Lock()
			//defer mutex.Unlock()
			t2 = fuzz_wordlist(words, url, t2)
		}()

	}
	wg.Wait()
	//remove the duplices from the slices
	uniqueMap := make(map[string]bool)
	uniqlists := []string{}
	for _, j := range foundlist {
		if !uniqueMap[j] {
			uniqueMap[j] = true
			uniqlists = append(uniqlists, j)
		}
	}

	//wire the ouput to a file
	writefile, _ := os.Create(outputfile)
	for _, i := range uniqlists {
		writefile.WriteString(i)
	}

	defer writefile.Close()
}
func fuzz_entensions(wordlist []string, url string, ext []string, uniqdisplay map[string]bool) map[string]bool {
	for _, word := range wordlist {
		for _, extension := range ext {
			comb := url + string(word) + "." + extension
			resp, err := http.Get(comb)
			if err != nil {
				fmt.Println("An Error Occured!!!")
			}
			if err == nil {
				if resp.StatusCode != 404 {
					body, _ := ioutil.ReadAll(resp.Body)
					combined := fmt.Sprintf("/%s : [length: %d size: %d Status: %d]\n", word+"."+extension, len(body), resp.ContentLength, resp.StatusCode)
					foundlist = append(foundlist, combined)
					for _, list := range foundlist {
						if !uniqdisplay[list] {
							uniqdisplay[list] = true
							color.Green("/%s : [length: %d size: %d Status: %d]\n", word+"."+extension, len(body), resp.ContentLength, resp.StatusCode)
						}
					}
				}
				resp.Body.Close()
			}

		}

	}
	return uniqdisplay
}

func fuzz_wordlist(wordlist []string, url string, uniqdisplay map[string]bool) map[string]bool {
	for _, word := range wordlist {
		comb := url + string(word)
		resp, err := http.Get(comb)
		if err != nil {
			fmt.Println("An Error Occured!!!")

		}
		if err == nil {
			if resp.StatusCode != 404 {
				body, _ := ioutil.ReadAll(resp.Body)
				combined := fmt.Sprintf("/%s : [length: %d size: %d Status: %d]\n", word, len(body), resp.ContentLength, resp.StatusCode)
				foundlist = append(foundlist, combined)
				for _, list := range foundlist {
					if !uniqdisplay[list] {
						uniqdisplay[list] = true
						color.Green("/%s : [length: %d size: %d Status: %d]\n", word, len(body), resp.ContentLength, resp.StatusCode)
					}
				}
			}
			resp.Body.Close()
		}

	}
	return uniqdisplay
}
