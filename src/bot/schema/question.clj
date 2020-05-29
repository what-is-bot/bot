(ns bot.schema.question
  (:require [bot.schema.answer :as schema.answer]
            [bot.schema.user :as schema.user]
            [schema.core :as s]))

(s/defschema Question {:id         s/Uuid
                       :created-at s/Inst
                       :asked-by   schema.user/User
                       :text       s/Str
                       :answers    [schema.answer/Answer]})
