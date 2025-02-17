// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
///////////////////////////////////////////////////////////////////////////////
#ifndef TINK_CORE_INTERNAL_PRIVATE_KEY_MANAGER_H_
#define TINK_CORE_INTERNAL_PRIVATE_KEY_MANAGER_H_

#include <memory>

#include "tink/core/internal_key_manager.h"
#include "tink/util/statusor.h"

namespace crypto {
namespace tink {

template <typename KeyProto, typename KeyFormatProto, typename PublicKeyProto,
          typename... Primitives>
class InternalPrivateKeyManager;

// An InternalPrivateKeyManager is an extension of InternalKeyManager. One
// should implement this in case there is a public key corresponding to the
// private key managed by this manager.
// Hence, in addition to the tasks a InternalKeyManager does, in order to
// implement a InternalPrivateKeyManager one needs to provide a function
// StatusOr<PublicKeyProto> GetPublicKey(const KeyProto& private_key) const = 0;
template <typename KeyProto, typename KeyFormatProto, typename PublicKeyProto,
          typename... Primitives>
class InternalPrivateKeyManager<KeyProto, KeyFormatProto, PublicKeyProto,
                                List<Primitives...>>
    : public InternalKeyManager<KeyProto, KeyFormatProto, List<Primitives...>> {
 public:
  explicit InternalPrivateKeyManager(
      std::unique_ptr<typename InternalKeyManager<KeyProto, KeyFormatProto,
                                                  List<Primitives...>>::
                          template PrimitiveFactory<Primitives>>... primitives)
      : InternalKeyManager<KeyProto, KeyFormatProto, List<Primitives...>>(
            std::move(primitives)...) {}

  virtual crypto::tink::util::StatusOr<PublicKeyProto> GetPublicKey(
      const KeyProto& private_key) const = 0;
};

}  // namespace tink
}  // namespace crypto

#endif  // TINK_CORE_INTERNAL_PRIVATE_KEY_MANAGER_H_
