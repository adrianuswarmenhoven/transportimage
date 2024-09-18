Read the [accompanying article](https://github.com/adrianuswarmenhoven/transportimage/blob/main/HideAnythingInEverything.md) first.

# Proof of Concept for transporting images through data and/or code 

Please take notice of the [License](LICENSE.md) before you download and/or use any code from this repo.

**Do not run any executable from this repo unless you know what you are doing!**

## Rationale

In cybersecurity, a lot of faith is put into so called 'techno-fixes': technological solutions to a security issue. Whilst this may be effective for some forms of malicious behavior (certain types of malware, phishing etc...), it can not cover the whole of 'data transference' between people.

However, there exists a school of thought that 'if you monitor everything, you will find everything'.

This, I think, is a dangerous and damaging way of thinking; the measures will break a lot of security and privacy for people and are of such a nature that these measures can not nor ever will be rolled back.

So, besides arguing the ethical side, I also want to demonstrate the technical futility of this way of thinking.

## About the methodology

I tried to find a good few examples that can clarify my point, whilst also trying to avoid helping people with less than good intentions.

To this end, the examples do work and can currently be used to avoid detection methods and circumvent moderation efforts, but they are cumbersome and inefficient. It should, however, be quite clear that if enough effort is put into this way of concealment, detection efforts will fail and things might become efficient really quick.

## What do the PoC's show?

I have generated an AI image of a kid holding up a thumb (and yes, the AI signature 6-fingers are there). The PoC tools will work with any image actually, but to demonstrate, I needed a copyright free image.

![AI generated image of Johnny Sixfingers](images/example_image.jpg "Johnny Sixfingers")

Each of the PoC's reads the image and then 'hides' or 'encodes' or 'steganographs' that image into HTML table, CSS, whitespace a spreadsheet, a binary or network traffic.

The result will not be flagged by client-side scanning, or by any AV. Moderators will see files of a different type and most likely will not deem it worthy of any further investigation.

Again, once an entity *knows* that this 'hiding' is being done, they can, relatively easy, find the content.

In some cases the recipient does not need to do anything special to view the image (CSS, HTML table, executable) and in some cases the recipient needs to do something to 'decode' or just to display the image (PCAP, Spreadsheet, Whitespace).

Each of the methods has it's own description.

## PoC's

- [Image to CSS](imageToCSS/)
- [Image to HTML Table](imageToTable/)
- [Image to spreadsheet](imageToSpreadsheet/)
- [Image to Executable](imageToExecutable/)
- [Image to PCAP](imageToPCAP/)
- [Image to whitespace](imageToWhitespace/)

## "Here's one I prepared earlier!"

Of course not everybody wants to go around compiling stuff and then running things etc... no, people are used to having things 'batteries included, ready to go!'.

So, I have supplied a folder named [output](output/) which has the results of the PoC's. Mind you, for anything but Linux I have been using cross compile, so I have no idea whether some things work. If something fails for your OS, please have a crack at compiling it from the PoC dir yourself first. 
