# Hide anything in everything
or, *'when you try to see everything, you will miss something'*

## There is no way

There is no way to monitor everything. Intrinsically, everybody knows this. However, the promise of a [Panopticon](https://en.wikipedia.org/wiki/Panopticon) still intrigues many a manager or policy maker in security.
This means every now and then (and in the last few years it was mostly the 'now') we get technical tools and/or policies that flirt with the panopticon idea, or that at the very least want try and get the benefits from a tightly monitored populace.

The truth is, that trying to monitor *everything* on the internet is like trying to monitor all of the atoms.

## There is too much

If you were to try and monitor all of the atoms, you would, sooner or later, reach the conclusion that 'To monitor something, let alone all of the atoms, I need something to monitor with, and that thing consists of atoms, so I can not monitor *all* of the atoms.'.
But even if you could manage to monitor it all, you would need atoms to record what the other atoms did etc. etc.
Eventually, you would reach the conclusion that it is wholly impossible to monitor all of the atoms. 
Well, then you only start monitoring meaningful atoms, right?
And this, this is when the problems start, because now there is a subjective measure introduced. And subjective means that *someone* decides about what is meaningful.
(This is the *context trap*. Atoms have different 'functions' when put into different configurations with other atoms (the configurations are colloquially known as 'molecules'))

Let's leave the atoms to themselves and let's get back to the internet. 

It will come as no surprise at all when I say that we have some analogy here, between the atoms and bits that are transmitted across the net.
There are simply too many bits to monitor in a meaningful way; Internet Exchanges pushing in excess of 15 Terabits at the moment of writing. That is, every second 15 x 10¹² bits are moved around.
And yes, this is not regular traffic, but it shows the immense throughput that gives us over 4500 Exabytes (1 Exabyte is 10¹² Megabytes) a year (again, at the time of writing, with [data from 2023](https://en.wikipedia.org/wiki/High_Bandwidth_Memory)).

That is a lot of bytes. And remember, if you want to just *read* all of that, it will take you a while to read it into memory, even with [HBM3](https://en.wikipedia.org/wiki/High_Bandwidth_Memory). Just reading all of that will take ages, let alone process it.

So, let's just monitor meaningful bytes, that should make the load considerably lighter, right?
Rrrrright.

## When viewed from this side...

Context and viewpoints matter. Imagine looking at a bowl from the top and looking at it from the bottom. One viewpoint will give you a concave surface, the other a convex.
The truth here is that there is one bowl and it has both a convex and a concave surface, and the point of view determines which you see.

This is what happens with data as well; all of the text, code, images, audio and video that gets send over the internet is sent as 0 and 1. 
The context of these zeroes and ones make that we can reconstruct, at the receiving side, what was sent.
And even without looking at sending and receiving, anything that you do with any digital device is just moving around a lot of 0's and 1's. The context then determines what effect these two numbers have: it can be a GIF or a PDF or a video...

The simple fact is that what is on your screen... it does not mean anything. It is us, humans, that *make* it mean *something*.
Our senses, together with our brain, decide that a group of color splotches in various shades, combined with some frequencies that are produced in a particular sequence, tell us we are watching other people make music (for instance). And unless a creature's senses have the same boundaries and limits as ours, together with a brain that has some things hardwired the same way we have, that creature is not going to experience the same thing that we are.

And neither is a computer. It does not experience things at all, it does what it is told to. Now, in order for it to assist us in thinking, reasoning, observing, we can 'teach' it, or 'train' it, to 'see' the same things we see. And we can train it very well.
But... it will still all be within the boundaries of what we normally experience.

Meaningful, that is from a shared and well-defined experience, that is.

## Fuzzy data

If you look at a lot of the digital artefacts that we peruse every day, you could notice that these artefacts leave quite some *'room for modification without distorting the original too much'*
Well, what that means is that I can make many literally different versions that people would absolutely not notice being different at all.
Basically, it shows us that there is a lot of 'usable' or 'malleable' data in a piece of content that does not change the intent, experience or effect of that content.
Now, if you add this piece of knowledge to the insane amount of content that is being moved around, you may get a feeling on how much 'fluff' space there is to play around with.

A more concrete explanation to follow. Think of a book. A book has one or more stories in it. These stories are made of letters, numbers and a selection of other characters. The sequence of these letters and their proximity to each other is important. 

```
Imeanthisisnotnicetoreadbutyoucanbecauseyoualreadyhavereadsomuchthatyouknowwordsandknowwhentheyend. 
```

```
Th isi sal otha r de rto re adbeca uset hewo rdbo unda rie sare unc lea r.
```

And that is why, in a book, with the stories, we stick to 'words have letters close to each other and between words we have some spacing so we see the groups of letters/words a lot easier. But...

```
This  is  totally  fine  to  read  even  though  there  are  two  spaces  between  each  word.
```

Hey now! That is a nice realization! We can modify the spaces in a story without making things too hard to read and without changing the content/meaning/story.

*As an excercise for the reader:* Think about the 0's and 1's again. How hard would it be to have some text using 1 space for 0 and two spaces for a 1? Would it show up when rendered as HTML? Are there different spaces with the same width but with a different meaning? 


## Steganography

Keep the insights from the previous paragraphs in mind, because we need to take a small detour into some theory of hiding things in things.

[Steganography](https://en.wikipedia.org/wiki/Steganography), which *"is the practice of representing information within another message or physical object, in such a manner that the presence of the concealed information would not be evident to an unsuspecting person's examination"* has been around for a long time (the term, from the Wikipedia article, at least since 1499, the practice, according to that same article, since 440 BC).

So, the concept of what is being done here is not new at all. In the Wikipedia article there are many examples which highlights what could be done.

The examples in this repo fall under some of those categories, but what I want to demonstrate is not steganography, it is not even that you could hide anything in just about everything.
I want to demonstrate that the *detection and enforcement* will always be a step behind.
This is just so I can point to this article when I talk about the protection of privacy versus draconian measures to find and/or catch a specific type of crime.

Let's get into the detection side. The [Steganalysis](https://en.wikipedia.org/wiki/Steganalysis) part.
From the Wikipedia article: *"Unlike cryptanalysis, in which intercepted data contains a message (though that message is encrypted), steganalysis generally starts with a pile of suspect data files, but little information about which of the files, if any, contain a payload. The steganalyst is usually something of a forensic statistician, and must start by reducing this set of data files (which is often quite large; in many cases, it may be the entire set of files on a computer) to the subset most likely to have been altered."*

And there we have it. The point I am going to make in this quite lengthy write-up for what basically is the battle of steganography versus steganalysis; '...starts with a pile of **suspect** data files...'.

There is so much data online that can be accessed by anyone at any time, that the first requirement, the logistics of identifying and then handling of suspect files, becomes the main problem.

A good Steganalyst can probably detect, in the case of the code in this repo quite easily (which is on purpose), the existence of hidden data.
However, that analyst would need to have selected that specific content type and those specific files prior to doing any analysis in the first place.

## But what does it mean?

As we saw from the previous paragraphs, it is, with some creativity, possible to hide anything in just about everything (and yes, I even, successfully, tried to hide things in the hidden things, like... steganography in an image which is hidden in a book).

We also saw that if there is a suspicion of Steganography, it is quite possible to detect it by various means.

And finally, we saw that the problem is not in the hiding versus detecting, but in the selection of where to check if something is hidden.

If you now imagine that anything you see online might contain a hidden message, you can start to get a feeling of the selection problem.
That selection problem is what I want to make you understand.

## The selection problem

It is currently technically impossible (and yes, this might be one of those articles that will age horribly, but the issue is now, at the time of writing).
If you just calculate the amount of bytes that need to go through the RAM and CPU, you will find that not only will it take an enormous amount of time, it will also be very costly.

The selection problem, or candidate problem, is one of the harder things in cybersecurity: in malware detection distinguishing between a benign and a malicious piece of content is hard and usually done with code that has been focused on the differences between the two.
Just getting good training data is a difficult excercise in itself.

The difficulty in selecting what content to analyse for steganography is off the charts though. Unless you already have some clues as to where to look, it means that you would have to analyze literally *everything* to just get the candidates. That would entail that you would know what that content 'normally' would be and how much of an aberration you would consider to be suspicious. If you have ever randomly browsed over [Github](github.com), you know that just defining the tolerances would be a monumental task.

## No Panopticon for you

What I am trying to get across is that (automated) monitoring of everything on every device, even if you outsource a lot of computing to end-user devices (or sources), there are ways around it, unless...

Unless...

Unless we start defining what is allowed to pass and what not. That would mean that governments would define rigid content format standards and that every single piece of content that you create will be statistically checked before it is allowed to be seen by others. You are basically only allowed to modify certain parts of the content in a way that is checked against what is 'normal'.

In security we have the concepts of 'default deny' and 'zero trust'. The internet is rather built on 'if the next machine in line accepts it, you are good'. Which essentially means 'if we can move these bytes across to another machine, we're good'. You can define your own protocol and use any port number (yes, even the IANA assigned ones), and you can send bytes in both directions.

And as long as the internet is allowed to innovate (although for some reason we are a bit hung up on using HTTP style protocols) and change, and as long as people are allowed to be creative, the whole 'catch them before the act' is just something that, as [Public Enemy once very eloquently stated: "It might feel good, it might sound a little somethin'
But damn the game if it ain't sayin' nothin'"](https://www.youtube.com/watch?v=7FmPskTljo0)

However, all the attempts to still panopticize the world, those cause more long-term harm than the short-term dopamine rush of catching the first few batches of criminals who will then pivot to using something else whilst the infrastructure for monitoring stays in place.

But that, and the whole issue of 'what is considered aberrant, malicious and who decides that' are matters for a different article.