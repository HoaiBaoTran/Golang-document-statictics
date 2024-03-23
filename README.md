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
Number of words: 147							
Number of characters: 778						
Average word length: 5.616822				
Execution Time 686.541µs		

file.txt 											

Number of lines: 125565				
Number of words: 3762705			
Number of characters: 20490810			
Average word length: 6.793051			
Execution Time 1.04579525s			

file_1MB.txt	 (basically file.txt * 5 times)			

Number of lines: 502265				
Number of words: 15050820				
Number of characters: 81963240				
Average word length: 6.793051					
Execution Time 3.939262209s	
	
——————————————————

With concurrency:	
Folder: StatisticsConcurrency
Folder: results - contain results
Command: go run main.go

requirement_en.txt									
Number of lines: 23								
Number of words: 147							
Number of characters: 778						
Average word length: 5.616822				
Execution Time 236.5µs

file.txt 											

Number of lines: 125565				
Number of words: 3762705			
Number of characters: 20490810			
Average word length: 6.819000		
Execution Time 264.295042ms

file_1MB.txt	 (basically file.txt * 5 times)			

Number of lines: 502265				
Number of words: 15050820				
Number of characters: 81963240				
Average word length: 6.819000	
Execution Time 979.6915ms

Method concurrency: Read all lines of file and split into chunks (default = 5)
Run concurrency statistics each chunk and merge them with fan in pattern

