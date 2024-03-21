# Golang-document-statictics

Statistics Document

Test 3 file: requirement_en.txt, file.txt and file_1MB.txt

Without concurrency: 
Folder: TextStatistics	
Folder: results - contain results
Command: go run TextStatistics.go				

Result:
requirement_en.txt									
Number of lines: 23								
Number of words: 150							
Number of characters: 817						
Average word length: 5.712963					
Execution Time 1.514375ms		

file.txt 											

Number of lines: 125565				
Number of words: 3762705			
Number of characters: 20490810			
Average word length: 6.772658			
Execution Time 2.456965333s			

file_1MB.txt	 (basically file.txt * 5 times)			

Number of lines: 502265				
Number of words: 15050820				
Number of characters: 81963240				
Average word length: 6.772658					
Execution Time 9.972920916s		
	
——————————————————

With concurrency:	
Folder: StatisticsConcurrency
Folder: results - contain results
Command: go run main.go

requirement_en.txt									
Number of lines: 23								
Number of words: 150							
Number of characters: 817						
Average word length: 5.712963					
Execution Time 854.666µs

file.txt 											

Number of lines: 125565				
Number of words: 3762705			
Number of characters: 20490810			
Average word length: 6.772658		
Execution Time 1.403031292s

file_1MB.txt	 (basically file.txt * 5 times)			

Number of lines: 502265				
Number of words: 15050820				
Number of characters: 81963240				
Average word length: 6.772658	
Execution Time 6.203735042s

Method concurrency: Read all lines of file and split into chunks (default = 5)
Run concurrency statistics each chunk and merge them with fan in pattern

