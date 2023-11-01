# gogrep
This is my rendition of the linux grep tool. Users input a search term for the root directory to search for within its files.
Individual workers concurrently search the files. They will return the file path to the found term, the line that include
the term, and the line number correlating to that search.
