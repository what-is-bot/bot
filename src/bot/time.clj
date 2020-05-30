(ns bot.time
  (:import [java.util Date]
           [java.time LocalDateTime ZoneOffset]))

(defn now [] (LocalDateTime/now))

(defn inst [] (Date.))

(defn time->inst [^LocalDateTime time]
  (-> time
      (.atZone (ZoneOffset/UTC))
      (.toInstant)
      Date/from))

(defn inst->time [^Date date]
  (-> date
      (.toInstant)
      (LocalDateTime/ofInstant ZoneOffset/UTC)))
