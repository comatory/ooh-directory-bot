package main

import (
	"testing"
)

func TestParsedHtmlWithSingleResult(t *testing.T) {
	results, _ := ParseResults(`
!DOCTYPE html
<html lang="en">
  <body>
	  <section class="grid">
		  <div class="grid__content">
			  <ol class="websites">
				  <li class="websites__item">
						<p class="website__intro">
							<a href="https://acme.xyz/">Acme site</a>
							<br>
							<q>Home page of Acme Inc.</q> By John Appleseed.
							<span class="text-nobr">
								<span class="text-emoji">
									<span aria-label="United States of America" title="United States of America">🇺🇸</span>
								</span>
								<a href="/blog/99999/">More info</a>
							</span>
						</p>
						<details class="website__details">
							<summary>
								Updated <time datetime="2022-11-21T19:04:00+00:00">a year ago</time>
							</summary>
							<div class="website__details__body">
								<figure class="figure figure--post">
									<figcaption>
										<a href="https://acme.xyz/2022/11/entry.html">Entry</a>
									</figcaption>
									<blockquote>
									  abcdqwerty
									</blockquote>
									<footer>
										<small>
										By John Appleseed, 3 words
										</small>
									</footer>
								</figure>
							</div> 
						</details>
					</li>
				</ol>
			</div>
		</section>
	</body>
</html>
`)

	if len(results) != 1 {
		t.Errorf("Expected to have \"%d\" results, got \"%d\"", 1, len(results))
	}

	result := results[0]

	if result.url != "https://acme.xyz/2022/11/entry.html" {
		t.Errorf("Expected to have result URL \"%s\", got \"%s\"", "https://acme.xyz/2022/11/entry.html", result.url)
	}

	if result.title != "Entry" {
		t.Errorf("Expected to have result title \"%s\", got \"%s\"", "Entry", result.title)
	}

	if result.authorName != "John Appleseed" {
		t.Errorf("Expected to have result author \"%s\", got \"%s\"", "John Appleseed", result.authorName)
	}

	if result.summary != "abcdqwerty" {
		t.Errorf("Expected to have result summary \"%s\", got \"%s\"", "abcdqwerty", result.summary)
	}
}

func TestParsedHtmlWithMultipleResults(t *testing.T) {
	results, _ := ParseResults(`
!DOCTYPE html
<html lang="en">
  <body>
	  <section class="grid">
		  <div class="grid__content">
			  <ol class="websites">
				  <li class="websites__item">
						<p class="website__intro">
							<a href="https://acme.xyz/">Acme site</a>
							<br>
							<q>Home page of Acme Inc.</q> By John Appleseed.
							<span class="text-nobr">
								<span class="text-emoji">
									<span aria-label="United States of America" title="United States of America">🇺🇸</span>
								</span>
								<a href="/blog/99999/">More info</a>
							</span>
						</p>
						<details class="website__details">
							<summary>
								Updated <time datetime="2022-11-21T19:04:00+00:00">a year ago</time>
							</summary>
							<div class="website__details__body">
								<figure class="figure figure--post">
									<figcaption>
										<a href="https://acme.xyz/2022/11/entry.html">Entry</a>
									</figcaption>
									<footer>
										<small>
										By John Appleseed, 3 words
										</small>
									</footer>
								</figure>
							</div> 
						</details>
					</li>
				  <li class="websites__item">
						<p class="website__intro">
							<a href="https://example.com/">Example</a>
							<br>
							<q>Example web site</q>.
							<span class="text-nobr">
								<span class="text-emoji">
									<span aria-label="United States of America" title="United States of America">🇺🇸</span>
								</span>
								<a href="/1x92">More info</a>
							</span>
						</p>
						<details class="website__details">
							<summary>
								Updated <time datetime="2023-09-12T03:00:00+00:00">3 months ago</time>
							</summary>
							<div class="website__details__body">
								<figure class="figure figure--post">
									<figcaption>
										<a href="https://acme.xyz/1x92">Example blog post</a>
									</figcaption>
									<footer>
									</footer>
								</figure>
							</div> 
						</details>
					</li>
				</ol>
			</div>
		</section>
	</body>
</html>
`)

	if len(results) != 2 {
		t.Errorf("Expected to have \"%d\" results, got \"%d\"", 2, len(results))
	}
}
