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
}
