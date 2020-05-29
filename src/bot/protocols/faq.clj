(ns bot.protocols.faq)

(defprotocol Faq
  (ask! [this question] "Ask a new question")
  (query [this & args] "Search for questions containing the given args")
  (answer! [this question answer] "Provide a new answer for the given question")
  (upvote! [this answer who] [this answer who delta] "Upvote an answer to increase its score. Can provide an optional delta (default is 1)")
  (downvote! [this answer who] [this answer who delta] "Downvote an answer to decrease its score. Can provide an optional delta (default is 1)"))

(def IFaq (:on-interface Faq))
