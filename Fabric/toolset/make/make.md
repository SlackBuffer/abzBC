# `make` 命令
<!-- http://www.ruanyifeng.com/blog/2015/02/make.html -->
- `make` 是一个根据指定的 Shell 命令进行构建的工具。规定要构建哪个文件，依赖哪些源文件，当那些文件有变动时，该如何重新构建它
	
    ```bash
    # 第一步，确认 b.txt 和 c.txt 必须已经存在
    # 第二步使用 cat 命令 将这个两个文件合并，输出为新文件
    a.txt: b.txt c.txt
        cat b.txt c.txt > a.txt
    

    make # 默认 Makefile 或 makefile
    make -f rules.txt
    make --file=rules.txt
    ```

    - 若运行 `make` 命令时未指定目标，默认会执行 `Makefile` 文件的第一个目标
# `Makefile` 文件的格式
- 规则（[rules](https://www.gnu.org/software/make/manual/html_node/Rules.html)）
	
    ```makefile
    <target> : <prerequisites> 
    [tab]  <commands>
    ```

    - 构建的"目标"是必需的，不可省略；"前置条件"和"命令"都是可选，但两者之中至少存在一个
    - 第二行必须由一个 tab 键起首
## `target`
- 一个目标（target）就构成一条规则
- 目标通常是文件名，指明 `make` 命令所要构建的对象
- 目标可以是一个文件名，也可以是**多个文件名**，之间用空格分隔
- 目标还可以是某个操作的名字，称为"伪目标"（phony target）
	
    ```makefile
    .PHONY: clean
    clean: 
        rm *.o temp
    ```

    - 不用“伪目标”的写法，若当前目录正好有一个 `clean` 文件，则 `make clean` 不会被执行，因为 `make` 发现 `clean` 文件已存在，无需重新构建
    - 声明 `clean` 是"伪目标"之后，`make` 就不会去检查是否存在一个叫做 `clean` 的文件，而是每次运行都执行对应的命令
    - > [Special Built-in Target Names](https://www.gnu.org/software/make/manual/html_node/Special-Targets.html#Special-Targets)
## `prerequisites`
- 前置条件通常是一组文件名，之间用空格分隔
- 它指定了"目标"是否重新构建的判断标准：只要有一个前置文件不存在，或者有过更新（前置文件的 `last-modification` 时间戳比目标的时间戳新），"目标"就需要重新构建
	
    ```makefile
    result.txt: source.txt
        cp source.txt result.txt
    # 若 source.txt 不存在，必须再写一条规则来生成 source.txt
    source.txt:
        echo "this is the source" > source.txt
    ```

    - `source.txt` 没有跟前置条件，意味着它跟其他文件都无关，只要这个文件还不存在，每次调用 `make source.txt`，它都会生成
- 生成多个文件的写法
	
    ```makefile
    source: file1 file2 file3

    # make source 等同于
    make file1
    make file2
    make file3
    ```

    - 此处 `source` 是一个伪目标，只有 3 个前置文件，没有对应任何命令，执行 `make source` 会一次性构建 file1, file2, file3
        - 此处 file 不一定是目标文件，也可以是伪目标
## `commands`
- 命令（commands）由一行或多行的Shell命令组成它，它的运行结果通常就是生成目标文件
- 每行命令之前必须有一个 tab 键，若想用其他键，可以用内置变量 `.RECIPEPREFIX` 声明
	
    ```makefile
    # > 代替 tab
    .RECIPEPREFIX = >
    all:
    > echo Hello, world
    ```

- 每行命令在一个单独的 shell 中执行，shell之间没有继承关系
	
    ```makefile
    var-lost:
        export foo=bar
        echo "foo=[$$foo]" # 读不到 foo 的值
    

    # Solutions
    var-kept:
        export foo=bar; echo "foo=[$$foo]"

    var-kept:
        export foo=bar; \
        echo "foo=[$$foo]"
    
    .ONESHELL:
    var-kept:
        export foo=bar; 
        echo "foo=[$$foo]"
    ```

## 语法
- 注释 `#`
- 正常情况下，`make` 会打印每条**命令**，然后再执行，称作回声（echoing）
	
    ```makefile
    test:
        # 这是测试
    
    test:
        @# 这是测试
    ```

    - 命令前加 `@` 可以关闭回声
    - 由于在构建过程中，需要了解当前在执行哪条命令，所以通常只在注释和纯显示的 `echo` 命令前面加上 `@`
- Makefile 的通配符与 Bash 一致，主要有星号（*）、问号（？）和 [...] 
- `make` 命令允许对文件名进行类似正则运算的匹配，主要用到的匹配符是 `%`
	
    ```makefile
    # 当前目录下有 f1.c 和 f2.c 2 个源码文件
    %.o: %.c
    # equivalent to
    f1.o: f1.c
    f2.o: f2.c

    # old-fashioned suffix Rules
    .c.o:
        $(CC) -c $(CFLAGS) $(CPPFLAGS) -o $@ $<
    ```

    - > https://www.gnu.org/software/make/manual/html_node/Pattern-Match.html#Pattern-Match
    - > https://www.gnu.org/software/make/manual/html_node/Suffix-Rules.html
- Makefile 允许使用 `=` 定义变量，调用时变量要放在 `$()` 中
	
    ```makefile
    txt = hello world
    test:
        @echo $(txt)
    ```

- 调用 Shell 变量，需要在美元符号前再加一个美元符号（因为 `make` 命令会对美元符号转义）
	
    ```makefile
    test:
        @echo $$HOME
    ```

- 4 个赋值运算符
	
    ```bash
    ## lazy set
    # Normal setting of a variable - values within it are recursively expanded 
    # when the variable is used, not when it's declared
    VARIABLE = value

    ## immediate set
    # Setting of a variable with simple expansion of the values inside
    # values within it are expanded at declaration time
    VARIABLE := value

    ## set if absent
    # Setting of a variable only if it doesn't have a value
    VARIABLE ?= value

    ## append
    # Appending the supplied value to the existing value 
    # or setting to that value if the variable didn't exist
    VARIABLE += value
    ```

    - `v1 = $(v2)` 中，变量 `v1` 的值是另一个变量 `v2`，`v1` 的值是在定义时扩展（静态扩展），还是在运行时扩展（动态扩展），结果可能会差异很大
    - https://www.gnu.org/software/make/manual/html_node/Setting.html
- `make` 命令提供一系列内置变量（Implicit Variables），主要是为了跨平台的兼容性
    - `$(CC)` 指向当前使用的编译器，`$(MAKE)` 指向当前使用的 `make` 工具
    - `CFLAGS`: Extra flags to give to the C compiler
    - > https://www.gnu.org/software/make/manual/html_node/Implicit-Variables.html
- `make` 命令还提供一些自动变量（*Automatic Variables*），它们的值与当前规则有关
    1. `$@` 指代当前目标，就是 `make` 命令当前构建的目标（`make foo` 的 `$@` 指代 `foo`）
    	
        ```makefile
        a.txt b.txt
            touch $@
        # equivalent to
        a.txt:
            touch a.txt
        b.txt: 
            touch b.txt
        ```
    
    2. `$<` 指代第一个前置条件（`t: p1 p2` 的 `$<` 指代 `p1`）
    3. `$?` - The names of **all the prerequisites** that are **newer** than the target, with spaces between them
    4. `$^` - The names of **all the prerequisites**, with spaces between them
    5. `$*` - The stem with which an implicit rule matches (If the target is `dir/a.foo.b` and the target pattern is `a.%.b` then the stem is `dir/foo`)
    6. `$(@D)`, `$(@F)` - 分别指向 `$@` 的目录名和文件名（`$@` 是 `src/input.c`，`$(@D)` 的值为 `src`，`$(@F)` 的值为 `input.c`）
    7. `$(<D)`, `$(<F)` - 分别指向 `$<` 的目录名和文件名
    - > https://www.gnu.org/software/make/manual/html_node/Automatic-Variables.html
    
    	
        ```makefile
        # [ ] -d dest?
        # 首先判断 dest 目录是否存在，不存在就新建
        # 将 src 目录下的 txt 文件，拷贝到 dest 目录下
        dest/&.txt: src/%.txt
            @[ -d dest ] || mkdir dest
            cp $< $@
        ```
    
- Makefile 使用 Bash 语法完成判断和循环
	
    ```makefile
    ifeq ($(CC), gcc)
        libs=$(libs_for_gcc)
    else
        libs=$(normal_libs)
    endif

    LIST = one two three
    all:
        for i in $(LIST); do \
            echo $$i; \
        done
    # equivalent to
    all:
        for i in one two three; do \
            echo $$i; \
        done
    ```

- 函数
	
    ```makefile
    $(function arguments)
    # or
    ${function arguments}
    ```

- 内置函数
    1. `shell`
    	
        ```makefile
        srcfiles := $(shell echo src/{00..99}.txt)
        ```
    
        - 用于执行 Shell 命令
    2. `wildcard`
        - Wildcard expansion happens automatically in rules. But wildcard expansion does not normally take place when a variable is set, or inside the arguments of a function. If you want to do wildcard expansion in such places, you need to use the `wildcard` function (`srcfiles := $(wildcard src/*.txt)`)
    3. `subst`
    	
        ```makefile
        $(subst ee,EE,feet on the street)
        # fEEt on the strEEt

        comma:= ,
        empty:=
        # space 变量用两个**空变量**作为标识符，中间是一个空格
        space:= $(empty) $(empty)
        foo:= a b c
        bar:= $(subst $(space),$(comma),$(foo))
        # bar is now `a,b,c'.
        ```
    
        - 文本替换 `$(subst from,to,text)`
    4. `patsubst`
    	
        ```makefile
        $(patsubst %.c,%.o,x.c.c bar.c)
        # x.c.o bar.o
        ```
    
    5. 替换后缀名
    	
        ```makefile
        # 将变量 OUTPUT 中的后缀名 .js 全部替换成 .min.js
        min: $(OUTPUT:.js=.min.js)
        ```
    
        - 变量名 + 冒号 + 后缀名替换规则；是 `patsubst` 函数的一种简写形式
    - > https://www.gnu.org/software/make/manual/html_node/Functions.html