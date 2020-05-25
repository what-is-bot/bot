(ns bot.system
  (:require [com.stuartsierra.component :as component]))

(def system (atom nil))

(defn bootstrap! [system-map]
  (->> system-map
       component/start-system
       (reset! system)))

(defn destroy! []
  (when-let [existing @system]
    (println "destroying system")
    (component/stop-system existing)
    (println "resetting system")
    (reset! system nil)
    (println "destroyed")))

(defn ensure-system-up! [system-map]
  (or @system (bootstrap! system-map)))
