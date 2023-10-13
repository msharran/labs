set files $HOME/sharran/**

for f in $files
    if test -n $f
        echo $f
    end
end
