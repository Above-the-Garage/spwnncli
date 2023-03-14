$dicts = @(23000,370000,49000)
$dicts |% {
    $f = "knownWords-" + $_ + ".txt"
    $w = cat $f
    $c = $w |% { $_.length }
    $a = $c | measure-object -average
    $f + " " + $a.Average
}