(ns bot.utils
  (:require [clojure.pprint :refer [pprint]]))

(def color-start "\033[7;37;44m\033[1m")
(def no-color "\033[0m")

(defn paint-text [text]
  (str color-start text no-color))

(defn assoc-some
  ([m k v]
   (if (some? v)
     (assoc m k v)
     m))
  ([m k v & kvs]
   {:pre [(even? (count kvs))]}
   (let [nm (assoc-some m k v)]
     (if (seq kvs)
       (let [[k v] (take 2 kvs)]
         (recur nm k v (nnext kvs)))
       nm))))

(defn conj-some
  ([coll v]
   (if (some? v)
     (conj coll v)
     coll))
  ([coll v & vs]
   (let [nc (conj-some coll v)]
     (if (seq vs)
       (recur nc (first vs) (next vs))
       nc))))

(defn- empty-or-multiple-of-3 [v]
  (boolean (or (not (seq v))
               (zero? (mod (count v) 3)))))

(defn assoc-pred
  ([m p k v]
   (if p
     (assoc m k v)
     m))
  ([m p k v & pkvs]
   {:pre [(empty-or-multiple-of-3 pkvs)]}
   (let [nm (assoc-pred m p k v)]
     (if (seq pkvs)
       (let [[p k v] (take 3 pkvs)]
         (recur nm p k v (next (nnext pkvs))))
       nm))))

(defn assoc-nil-key
  ([m k v]
   (assoc-pred m (nil? (get m k)) k v))
  ([m k v & kvs]
   (let [nm (assoc-nil-key m k v)]
     (if (seq kvs)
       (let [[k v] (take 2 kvs)]
         (recur nm k v (nnext kvs)))
       nm))))

(defn pprint-to-str [v]
  (let [printed (with-out-str (pprint v))]
    (subs printed 0 (dec (count printed)))))

(defn debug [v]
  `(let [v# ~v]
    (println
    (paint-text "debug=> ")
    (pprint-to-str ~v)
    (paint-text " <="))
    v#))

(defn debugd [v]
  `(let [v# ~v]
     (println (str (paint-text "debugd")
                   " "
                   '~v
                   " "
                   (paint-text "=>")
                   \newline
                   (pprint-to-str v#)
                   " "
                   (paint-text "<=")))
     v#))
