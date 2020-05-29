(ns bot.schema.answer
  (:require [schema.core :as s]
            [bot.schema.user :as schema.user]))

(s/defschema Answer {:id          s/Uuid
                     :created-at  s/Inst
                     :answered-by schema.user/User
                     :text        s/Str
                     :score       s/Int})
