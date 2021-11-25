## IP Counter Demo

This is a demo of maintaining counts of IP hits and a top 100 list of the IPs with the most hits.

### How to Run

````
$ go build
# macOS or Linux:
$ ./ip-counter
# Windows:
$ .\ip-counter.exe
````

### Data Structures

The normal IP hits are maintained in a simple hash map of IP (as int) to a 64 bit count. This uses as minimal space as possible and is quick to update.

The top 100 list is a maintained as a sorted doubly-linked list using the built-in Go data structure. This is more or less a max heap.

The process of keeping the list sorted is fairly complicated but generally only involves inserting a new value or moving an existing value. Since the counts are only changing by 1 each time, any movement will usually be short. Also the list only needs to be updated when an IP's count is more than or equal to the lowest value currently in the list. The runtime is generally O(1).

The end result is a top 100 list runtime which is so small on my Ryzen desktop that it shows up as zero.

Also on average the runtime of `RequestHandled` is generally on the order of 400 nanoseconds (around half a millisecond).

### Main function

This has a `main()` which acts as a bit of a demo. Here is an example run:

````
Took 193.1762ms to generate 500000 counts for a total of 500 IPs
Average time to handle each IP is 386ns
Took 0s to get top 100
1. IP: 0.0.1.58, Count: 982
2. IP: 0.0.0.244, Count: 982
3. IP: 0.0.1.120, Count: 981
4. IP: 0.0.0.137, Count: 980
5. IP: 0.0.1.218, Count: 980
6. IP: 0.0.1.191, Count: 979
7. IP: 0.0.0.219, Count: 978
8. IP: 0.0.0.255, Count: 978
9. IP: 0.0.0.161, Count: 977
10. IP: 0.0.1.64, Count: 977
11. IP: 0.0.0.216, Count: 976
12. IP: 0.0.1.112, Count: 976
13. IP: 0.0.0.68, Count: 975
14. IP: 0.0.1.50, Count: 975
15. IP: 0.0.0.13, Count: 974
16. IP: 0.0.0.6, Count: 974
17. IP: 0.0.1.209, Count: 974
18. IP: 0.0.1.162, Count: 973
19. IP: 0.0.0.69, Count: 973
20. IP: 0.0.0.131, Count: 973
21. IP: 0.0.0.57, Count: 973
22. IP: 0.0.0.222, Count: 972
23. IP: 0.0.0.79, Count: 972
24. IP: 0.0.0.157, Count: 971
25. IP: 0.0.0.72, Count: 970
26. IP: 0.0.1.211, Count: 970
27. IP: 0.0.0.232, Count: 970
28. IP: 0.0.1.178, Count: 969
29. IP: 0.0.1.182, Count: 969
30. IP: 0.0.1.99, Count: 969
31. IP: 0.0.1.234, Count: 969
32. IP: 0.0.1.88, Count: 968
33. IP: 0.0.1.35, Count: 967
34. IP: 0.0.0.192, Count: 966
35. IP: 0.0.1.196, Count: 966
36. IP: 0.0.1.11, Count: 966
37. IP: 0.0.1.103, Count: 965
38. IP: 0.0.1.124, Count: 965
39. IP: 0.0.1.95, Count: 965
40. IP: 0.0.1.54, Count: 965
41. IP: 0.0.1.150, Count: 964
42. IP: 0.0.0.14, Count: 964
43. IP: 0.0.1.111, Count: 964
44. IP: 0.0.1.188, Count: 963
45. IP: 0.0.1.242, Count: 962
46. IP: 0.0.1.189, Count: 962
47. IP: 0.0.0.170, Count: 962
48. IP: 0.0.0.155, Count: 962
49. IP: 0.0.1.51, Count: 962
50. IP: 0.0.1.151, Count: 962
51. IP: 0.0.0.71, Count: 962
52. IP: 0.0.1.225, Count: 961
53. IP: 0.0.0.201, Count: 961
54. IP: 0.0.0.172, Count: 960
55. IP: 0.0.1.134, Count: 960
56. IP: 0.0.0.242, Count: 960
57. IP: 0.0.1.16, Count: 960
58. IP: 0.0.1.123, Count: 958
59. IP: 0.0.0.138, Count: 958
60. IP: 0.0.0.193, Count: 958
61. IP: 0.0.0.37, Count: 958
62. IP: 0.0.0.225, Count: 958
63. IP: 0.0.0.174, Count: 957
64. IP: 0.0.1.230, Count: 957
65. IP: 0.0.0.10, Count: 955
66. IP: 0.0.1.159, Count: 954
67. IP: 0.0.0.38, Count: 954
68. IP: 0.0.1.175, Count: 954
69. IP: 0.0.1.174, Count: 954
70. IP: 0.0.0.149, Count: 952
71. IP: 0.0.1.135, Count: 951
72. IP: 0.0.1.73, Count: 950
73. IP: 0.0.0.76, Count: 950
74. IP: 0.0.0.210, Count: 949
75. IP: 0.0.0.206, Count: 949
76. IP: 0.0.0.104, Count: 949
77. IP: 0.0.1.110, Count: 949
78. IP: 0.0.1.37, Count: 948
79. IP: 0.0.0.21, Count: 948
80. IP: 0.0.1.201, Count: 947
81. IP: 0.0.0.51, Count: 947
82. IP: 0.0.1.92, Count: 947
83. IP: 0.0.0.12, Count: 947
84. IP: 0.0.0.245, Count: 946
85. IP: 0.0.0.160, Count: 944
86. IP: 0.0.1.136, Count: 943
87. IP: 0.0.1.237, Count: 942
88. IP: 0.0.1.106, Count: 941
89. IP: 0.0.1.20, Count: 940
90. IP: 0.0.1.13, Count: 940
91. IP: 0.0.1.78, Count: 939
92. IP: 0.0.1.227, Count: 939
93. IP: 0.0.0.4, Count: 937
94. IP: 0.0.1.53, Count: 936
95. IP: 0.0.1.80, Count: 935
96. IP: 0.0.0.198, Count: 935
97. IP: 0.0.1.76, Count: 933
98. IP: 0.0.1.114, Count: 929
99. IP: 0.0.0.90, Count: 925
100. IP: 0.0.0.28, Count: 918
````