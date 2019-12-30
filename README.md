# cover_filter

		    ______                            ______ _  __ __			
		   / ____/____  _   __ ___   _____   / ____/(_)/ // /_ ___   _____	
		  / /    / __ \| | / // _ \ / ___/  / /_   / // // __// _ \ / ___/	
		 / /___ / /_/ /| |/ //  __// /     / __/  / // // /_ /  __// /	
		 \____/ \____/ |___/ \___//_/     /_/    /_//_/ \__/ \___//_/		
     
coverage_filter ———— 覆盖率报告修改工具，根据标记过滤代码块（最小单位：line）

  -count int
  
        指定注释标记块被统计次数，负数表示不统计改代码块 (default 1)
        
  -file string
  
        通过go test生成覆盖率的源文件文件 (default "result.out")
        
  -httptest.serve string
  
        if non-empty, httptest.NewServer serves on this address and blocks
        
  -marker string
  
        设置过滤标记，可设置多个标签，以逗号(,)分隔 (default "no-cover")
        
  -out string
  
        设置输出的文件名 (default "result.out.coverfilter")
        
