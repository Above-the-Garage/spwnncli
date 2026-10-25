# Dictionaries

This file records what is known about each word list in spwnncli: where it came
from, when it was obtained, its format and how it relates to the other files.

* **Fact** means verified: git history, a byte-identical upstream file, or a README/commit message.
* **Best guess** means an inference; the reasoning is given.

Git dates are first-commit dates in this repo ("added to repo"). File mtimes are
not evidence (a clone sets them all to the same date).

spwnn's `ReadDictionary` reads one word per line and maps anything outside a-z
(and `_`) to a single "other" letter, so apostrophes, hyphens, capitals and
accented letters all become the same symbol. Lists that are lowercase a-z only
are the cleanest input.

## Summary

| File | Lines | Bytes | Format | Same as |
|---|---|---|---|---|
| knownWords.txt (default) | 49,167 | 448,261 | lowercase, LF, 3,655 lines with `'` (possessives) | knownWords-49000.txt |
| knownWords-49000.txt | 49,167 | 448,261 | as above | knownWords.txt |
| knownWords-23000.txt | 23,115 | 200,777 | mixed case, LF, 307 with `'`, 22 with Latin-1 accented bytes | american.sml+.mwl minus affix flags |
| american.sml+.mwl | 23,115 | 269,976 | ispell munched list (`word/FLAGS`), mixed case, Latin-1 | |
| knownWords-370000.txt | 370,096 | 3,864,619 | lowercase a-z only, LF, no trailing newline | (dwyl words_alpha minus 15 concatenations, see below) |
| dwyl-english-words-words_alpha.txt | 370,103 | 3,864,812 | as above | |
| knownWords-scowl60.txt | 87,780 | 815,351 | lowercase a-z only, LF, sorted, unique | |
| knownWords-scowl70.txt | 127,770 | 1,220,221 | lowercase a-z only, LF, sorted, unique | |
| knownWords-shakespeare.txt | 12,878 | 104,136 | lowercase, LF; 755 entries with non a-z (hyphens such as `a-bed`, `&c`, `3`) | tinyshakespeare-words-sorted-cleaned-uniq.txt |
| tinyshakespeare.txt | 40,000 | 1,115,394 | original text, not a word list | |
| tinyshakespeare-words.txt | 209,892 | 1,115,394 | one token per line, punctuation and case kept | |
| tinyshakespeare-words-sorted.txt | 209,892 | 1,115,394 | sorted tokens | |
| tinyshakespeare-words-sorted-cleaned.txt | 209,892 | 1,062,006 | lowercased, punctuation stripped, blanks and duplicates kept | |
| tinyshakespeare-words-sorted-cleaned-uniq.txt | 12,878 | 104,136 | unique cleaned words | knownWords-shakespeare.txt |
| spwnnErrors-2022-11-22.txt | 1,041 | 152,522 | CRLF; a captured console log, not a word list | |
| contentedTest.txt | 56 | 1,243 | captured spwnncli output, not a word list | |

Identical pairs were verified with `cmp`.

## knownWords.txt / knownWords-49000.txt (the ~49k benchmark list)

* **Origin (Best guess, from Stephen):** obtained from a Unix system, possibly a
  DECSYSTEM-20, in the 1980s when this project started; alternatively from when
  the first C++ version was written. Pending a lookup in Stephen's archives.
* **Fact:** the spwnn README calls it "ispell.words - about 49,000 words (of
  somewhat unknown origin, circa 2000)". The README's history section says a C
  version was written "around 2000".
* **Fact:** commit 84b97fa (2022-11-22) "ispell_words for testing" added
  `ispell_words.txt`, byte-identical to the original knownWords.txt; it was
  deleted in 86955f2 the same day.
* **Obtained:** unknown (1980s or the first C++ version), pending archive lookup; earliest evidence in git: 2020-03-15. **Added to repo
  (Fact):** 2020-03-15, commit 4658ad5 "first commit".
* **History (Fact):** knownWords-49000.txt was split out in 8626ba0 (2022-11-22,
  "Make knownWords versions obvious"). cbe7f63 (2022-11-22) removed the typo
  `procotols` from both files. knownWords.txt was swapped to other lists and
  back several times (5cbf831, ff1c555, 7744c82, a357d8e, fd0619a); since
  4daa00b (2023-08-11, "Back to 49000 sized default dictionary.") it is
  byte-identical to knownWords-49000.txt.
* **Evidence (Fact, measured 2026-10-24):** all lowercase, including proper
  names (`aarhus`, `aaron`); 3,655 possessive entries (`abbey's`). 19,277 of the
  23,113 lowercased american.sml+.mwl words are in it. 42,440 of its words are
  in SCOWL size 50 (en_US, lowercased) and 6,727 are not.
* **Best guess (own):** the possessives and the overlap with the ispell list fit
  an ispell-family word list (the 2022 commit message calls it ispell_words).
  This doesn't settle whether it came from a 1980s Unix machine.
* Don't change it; it is the benchmark list.

## american.sml+.mwl and knownWords-23000.txt

* **Origin (Fact):** commit 5be30d2 (2022-12-31) "ispell small american
  dictionary"; the spwnn README lists it as "about 23,000 words (from ispell,
  2022)".
* **Best guess:** the ispell distribution's `american.sml+` list in munched
  form (`.mwl`), given the name and the `/S`, `/DGRS` affix flags.
* **Obtained:** 2022 (README); **added to repo (Fact):** 2022-12-31.
* **knownWords-23000.txt (Fact):** added in 5cbf831 (2022-12-31) as a copy of
  american.sml+.mwl; 7a2dbe0 ("Remove suffixes") stripped the `/FLAGS`. It is
  exactly `sed 's|/.*||' american.sml+.mwl` (verified). The affixes were removed,
  not expanded, so plurals and other inflections are missing. It keeps capitals,
  apostrophes and 22 Latin-1 accented bytes (for example `attaché` as byte 0xE9).

## dwyl-english-words-words_alpha.txt / knownWords-370000.txt

* **Origin (Fact):** `words_alpha.txt` from https://github.com/dwyl/english-words
  (commit f47fb06, 2022-11-22, "Big modern dictionary"). dwyl's README says the
  words came from an infochimps "350,000 simple English words" Excel file and
  that "Copyright still belongs to them". The spwnn README says "from aspell",
  which doesn't match dwyl's own account.
* **Added to repo (Fact):** 2022-11-22. knownWords-370000.txt was split out in
  8626ba0 the same day; the two were byte-identical until 2026-10-24.
* **Edited (Fact):** on 2026-10-24, 15 words in knownWords-370000.txt that were
  two words run together (aeonicaeonist, aliyahaliyahs, brigantinebrigantines,
  coevolvedcoevolves, demarcatordemarcators, epigonousepigons, gynecicgynecidal,
  kestrelkestrels, nabobrynabobs, outwardsoutwarred, peroxidicperoxiding,
  preinferredpreinferring, stegnosisstegnotic, sturdiersturdies,
  subsultorysubsultus) were split back into their halves. Halves already in the
  list were not re-added; the 8 new words are aliyahs, coevolved, coevolves,
  epigons, outwarred, preinferred, preinferring and sturdies. Many look like a
  dictionary page's two guide words joined. dwyl-english-words-words_alpha.txt
  is left unchanged as the original.
* **Snapshot date (Best guess):** compared with upstream on 2026-10-24, this copy
  still has `acceleratorh` (removed upstream 2020-08-24), has `cryptocurrency`
  (added 2020-04-10), and lacks `nerdy`/`abled` (added April 2022) and `cyber`
  (2025). It also still has `unneccessary`, which upstream removed in 2024. So
  the content matches upstream between April and August 2020, although it was
  added here in 2022 (possibly from an older copy or mirror). Upstream uses CRLF;
  this copy is LF with no trailing newline.
* **Quality (Fact, measured 2026-10-24):** it contains typos such as
  `accomodate`, `enviroment`, `supprise`, `wierd`, `untill` (and the rare real
  word `calender`). 258,305 of its words (70%) are not in SCOWL size 70; only
  8,080 are missing from SCOWL size 95 (en_US), mostly British spellings. So dwyl
  is roughly SCOWL 95, and SCOWL 95 itself contains those same typos (all from
  `english-words.95`).

## knownWords-scowl60.txt / knownWords-scowl70.txt

* **Origin (Fact):** SCOWL (Spell Checker Oriented Word Lists) release
  2020.12.07 by Kevin Atkinson, `scowl-2020.12.07.tar.gz` from
  https://sourceforge.net/projects/wordlist/files/SCOWL/2020.12.07/ (project
  page http://wordlist.aspell.net/).
* **Obtained (Fact):** 2026-10-24.
* **Recipe (Fact):** inside the unpacked release:

  ```sh
  (./mk-list english american 60; cat final/special-hacker.50) \
    | iconv -f ISO-8859-1 -t UTF-8 | tr 'A-Z' 'a-z' \
    | LC_ALL=C grep -E '^[a-z]+$' | LC_ALL=C sort -u > knownWords-scowl60.txt
  ```

  The same with `70`. This is US spelling with no variants (the en_US setup).
  It includes proper names (lowercased) and drops anything with an apostrophe,
  hyphen, digit or accent, so possessives and contractions are gone.
* **Counts (Fact):** size 60: 87,780 words; size 70: 127,770 words.
* **Test findings (Fact, `spwnncli -test`, 2026-10-24, M4 MacBook Air):**
  * Size 60: capped 1m24s, 4 ambiguous lines (2 pairs: `beavered`/`bereaved`,
    `indented`/`intended`), 2/20 misspellings failed (`seperate`, `enviroment`);
    original 1m28s, 36 ambiguous, the same 2 failures.
  * Size 70: capped 3m22s, 16 ambiguous, 3/20 failed (adds `calender`); original
    3m32s, 74 ambiguous, the same 3 failures.
  * None of the 20 test typos is in either list except `calender`, which is in
    size 70 (a real but rare word, a cloth-pressing machine), so `calender`
    corrects to itself there.
  * Size 70 contains `iaccoca` alongside `iacocca`, a misspelled proper name.
* **License (Fact):** see `SCOWL-Copyright.txt` (the release's `Copyright` file,
  included verbatim). It is permissive (MIT/BSD-like): use, copy, modify,
  distribute and sell are allowed, provided the copyright notice and permission
  notice appear in all copies and in supporting documentation. Most source lists
  up to size 70 are public domain (Moby, 12Dicts, Kelk's UK list); the file also
  carries WordNet 1.6, VarCon and ispell (Geoff Kuenning) notices that ask for
  the notice to be kept. Keep `SCOWL-Copyright.txt` with these files.

## tinyshakespeare*.txt and knownWords-shakespeare.txt

* **Origin (Fact):** tinyshakespeare.txt is byte-identical (MD5
  6fb458f1232090904fb40fe944165e91, 1,115,394 bytes) to Andrej Karpathy's
  char-rnn dataset `data/tinyshakespeare/input.txt`
  (https://github.com/karpathy/char-rnn).
* **Added to repo (Fact):** 2023-01-23, commit a357d8e "Add "tiny Shakespeare"
  as a dictionary." Obtained: unknown, earliest evidence 2023-01-23.
* **Derivation:** tinyshakespeare-words.txt is exactly tinyshakespeare.txt with
  spaces turned into newlines (**Fact**, verified). The `-sorted`, `-cleaned`
  and `-uniq` steps (sort; lowercase and strip punctuation; unique) are a **Best
  guess** from the names and contents; a plain `LC_ALL=C sort` doesn't reproduce
  the byte order exactly. knownWords-shakespeare.txt is byte-identical to
  tinyshakespeare-words-sorted-cleaned-uniq.txt (Fact). It contains Elizabethan
  spellings (`abandond`) and hyphenated compounds.

## spwnnErrors-2022-11-22.txt (not a dictionary)

* **Fact:** added in defbf6a (2022-11-22) "Add errors (mostly) caused by
  duplicate pairs in words for both dictionaries" and edited through
  2022-12-14. A CRLF console capture of `spwnncli -g` runs on an EC2 x1.32xlarge
  (Intel Haswell) with Go 1.18, against the 49,167-word and 370,104-word lists.

## contentedTest.txt (not a dictionary)

* **Fact:** added in 3fcd5ef (2022-12-24). Captured spwnncli output for the
  input `contented`, showing ties with `consented`, `contended`, `contested` and
  `concentrated`.

## Other

* AverageWordLen.ps1 (Fact, 6e85396, 2023-03-13) prints the average word
  length of knownWords-23000/-370000/-49000. It doesn't say where they came from.
