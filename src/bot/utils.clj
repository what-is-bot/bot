(ns bot.utils)

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
