(ns bot.schema.answer
  (:require [schema.core :as s]
            [bot.schema.user :as schema.user]))

(s/defschema Answer {(s/optional-key :id)          s/Uuid
                     (s/optional-key :created-at)  s/Inst
                     (s/optional-key :answered-by) schema.user/User
                     :text                         s/Str
                     (s/optional-key :score)       s/Int})
