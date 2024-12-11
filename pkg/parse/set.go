package parse

type set struct{
	chr map[byte]bool
}

func (s set) String () string {
    res := "["
    for k := range s.chr {
        res += string(k) 
    }
    res += "]"
	return res 
}