local func_with_opts = function(opts)
    if opts == nil then
        print("opts is nil")
        return
    end
    print("opts.a: ", opts.a)
    print("opts.b: ", opts.b)
end

print("Calling func_with_opts with no arguments")
func_with_opts()

print("Calling func_with_opts with table opts")
func_with_opts({ a = 1, b = 2 })

print("Calling func_with_opts with table opts with one key")
func_with_opts({ a = 1 })

print("Calling func_with_opts with table opts without param")
func_with_opts { a = 1, b = 2 }
