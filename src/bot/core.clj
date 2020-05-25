(ns bot.core
  (:gen-class)
  (:require [bot.components :as components]))

(defn -main [&]
  (components/create-and-start-system!))
