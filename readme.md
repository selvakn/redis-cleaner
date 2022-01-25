## redis-cleaner

### Usage
<pre>
./redis-cleaner:
  -P string
    	redis connection password (optional)
  -c string
    	redis connection string (default "localhost:6379")
  -key-prefix string
    	key prefix to search (if not mention, then all keys)
  -retain-days-max int
    	keys to be retained with ttl more than this value (default 450)
  -retain-days-min int
    	keys to be retained with ttl less than this value (default 15)
</pre>