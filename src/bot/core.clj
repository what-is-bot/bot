(ns bot.core
  (:gen-class)
  (:require [bot.components :as components]))

(defn -main [& _]
  (components/create-and-start-system!))
