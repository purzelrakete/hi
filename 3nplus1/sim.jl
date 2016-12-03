# try to show p2^n=q2^m
function smash()
  p = rand(1:2:100) # p is odd
  q = rand(1:2:100) # q is odd
  m = rand(1:100)
  n = rand(1:100)

  if p == q
    return
  end

  if p*2^n == q*2^m
    error("$p*2^$n == $q*2^$m")
  end
end

# ONE BILLIONNN, takes about 5 minutes.
for x in 1:1_000_000_000
  smash()
end
