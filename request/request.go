package request

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/DonggyuLim/Alliance-Rank/utils"
)

type AccountsReq struct {
	Account    []AccountReq `json:"accounts"`
	Pagination Pagination   `json:"pagination"`
}

type Pagination struct {
	NextKey string `json:"next_key"`
	Total   string `json:"total"`
}
type AccountReq struct {
	Type    string `json:"@type"`
	Address string `json:"address"`
	PubKey  struct {
		Type string `json:"@type"`
		Key  string `json:"key"`
	} `json:"pub_key"`
	AccountNumber string `json:"account_number"`
	Sequence      string `json:"sequence"`
}

type BalanceReq struct {
	Balances   []Balance `json:"balances"`
	Pagination struct {
		NextKey interface{} `json:"next_key"`
		Total   string      `json:"total"`
	} `json:"pagination"`
}

type DelegationRequest struct {
	Deligations []DelegationData `json:"delegations"`
	Pagination  interface{}      `json:"pagination,omitempty"`
}

type DelegationData struct {
	Delegation struct {
		DelegatorAddress      string          `json:"delegator_address"`
		ValidatorAddress      string          `json:"validator_address"`
		Denom                 string          `json:"denom"`
		Shares                string          `json:"shares"`
		RewardHistory         []RewardHistory `json:"reward_history,omitempty"`
		LastRewardClaimHeight string          `json:"last_reward_claim_height"`
	} `json:"delegation"`
	Balance Balance `json:"balance"`
}

func (d *DelegationData) EncodeByte(b []byte) {
	encoder := gob.NewDecoder(bytes.NewReader(b))
	utils.PanicError(encoder.Decode(&d))
}

type RewardHistory struct {
	Denom string `json:"denom"`
	Index string `json:"index"`
}
type Reward struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

/*
	{
	    "rewards": [
	        {
	            "denom": "ibc/D7AA592A1C1C00FE7C9E15F4BB7ADB4B779627DD3FBB3C877CD4DB27F56E35B4",
	            "amount": "35744"
	        },
	        {
	            "denom": "uatr",
	            "amount": "9231927"
	        }
	    ]
	}
*/
type RewardRequest struct {
	Rewards []Reward `json:"rewards,omitempty"`
}
type Balance struct {
	Denom  string `bson:"denom" json:"denom"`
	Amount string `bson:"amount" json:"amount"`
}

type LastBlock struct {
	BlockID struct {
		Hash          string `json:"hash,omitempty"`
		PartSetHeader struct {
			Total int    `json:"total,omitempty"`
			Hash  string `json:"hash,omitempty"`
		} `json:"part_set_header,omitempty"`
	} `json:"block_id,omitempty"`
	Block struct {
		Header struct {
			Version struct {
				Block string `json:"block,omitempty"`
				App   string `json:"app,omitempty"`
			} `json:"version,omitempty"`
			ChainID     string    `json:"chain_id,omitempty"`
			Height      string    `json:"height,omitempty"`
			Time        time.Time `json:"time,omitempty"`
			LastBlockID struct {
				Hash          string `json:"hash,omitempty"`
				PartSetHeader struct {
					Total int    `json:"total,omitempty"`
					Hash  string `json:"hash,omitempty"`
				} `json:"part_set_header,omitempty"`
			} `json:"last_block_id,omitempty"`
			LastCommitHash     string `json:"last_commit_hash,omitempty"`
			DataHash           string `json:"data_hash,omitempty"`
			ValidatorsHash     string `json:"validators_hash,omitempty"`
			NextValidatorsHash string `json:"next_validators_hash,omitempty"`
			ConsensusHash      string `json:"consensus_hash,omitempty"`
			AppHash            string `json:"app_hash,omitempty"`
			LastResultsHash    string `json:"last_results_hash,omitempty"`
			EvidenceHash       string `json:"evidence_hash,omitempty"`
			ProposerAddress    string `json:"proposer_address,omitempty"`
		} `json:"header,omitempty"`
		Data struct {
			Txs []string `json:"txs,omitempty"`
		} `json:"data,omitempty"`
		Evidence struct {
			Evidence []struct {
				DuplicateVoteEvidence struct {
					VoteA struct {
						Type    string `json:"type,omitempty"`
						Height  string `json:"height,omitempty"`
						Round   int    `json:"round,omitempty"`
						BlockID struct {
							Hash          string `json:"hash,omitempty"`
							PartSetHeader struct {
								Total int    `json:"total,omitempty"`
								Hash  string `json:"hash,omitempty"`
							} `json:"part_set_header,omitempty"`
						} `json:"block_id,omitempty"`
						Timestamp        time.Time `json:"timestamp,omitempty"`
						ValidatorAddress string    `json:"validator_address,omitempty"`
						ValidatorIndex   int       `json:"validator_index,omitempty"`
						Signature        string    `json:"signature,omitempty"`
					} `json:"vote_a,omitempty"`
					VoteB struct {
						Type    string `json:"type,omitempty"`
						Height  string `json:"height,omitempty"`
						Round   int    `json:"round,omitempty"`
						BlockID struct {
							Hash          string `json:"hash,omitempty"`
							PartSetHeader struct {
								Total int    `json:"total,omitempty"`
								Hash  string `json:"hash,omitempty"`
							} `json:"part_set_header,omitempty"`
						} `json:"block_id,omitempty"`
						Timestamp        time.Time `json:"timestamp,omitempty"`
						ValidatorAddress string    `json:"validator_address,omitempty"`
						ValidatorIndex   int       `json:"validator_index,omitempty"`
						Signature        string    `json:"signature,omitempty"`
					} `json:"vote_b,omitempty"`
					TotalVotingPower string    `json:"total_voting_power,omitempty"`
					ValidatorPower   string    `json:"validator_power,omitempty"`
					Timestamp        time.Time `json:"timestamp,omitempty"`
				} `json:"duplicate_vote_evidence,omitempty"`
				LightClientAttackEvidence struct {
					ConflictingBlock struct {
						SignedHeader struct {
							Header struct {
								Version struct {
									Block string `json:"block,omitempty"`
									App   string `json:"app,omitempty"`
								} `json:"version,omitempty"`
								ChainID     string    `json:"chain_id,omitempty"`
								Height      string    `json:"height,omitempty"`
								Time        time.Time `json:"time,omitempty"`
								LastBlockID struct {
									Hash          string `json:"hash,omitempty"`
									PartSetHeader struct {
										Total int    `json:"total,omitempty"`
										Hash  string `json:"hash,omitempty"`
									} `json:"part_set_header,omitempty"`
								} `json:"last_block_id,omitempty"`
								LastCommitHash     string `json:"last_commit_hash,omitempty"`
								DataHash           string `json:"data_hash,omitempty"`
								ValidatorsHash     string `json:"validators_hash,omitempty"`
								NextValidatorsHash string `json:"next_validators_hash,omitempty"`
								ConsensusHash      string `json:"consensus_hash,omitempty"`
								AppHash            string `json:"app_hash,omitempty"`
								LastResultsHash    string `json:"last_results_hash,omitempty"`
								EvidenceHash       string `json:"evidence_hash,omitempty"`
								ProposerAddress    string `json:"proposer_address,omitempty"`
							} `json:"header,omitempty"`
							Commit struct {
								Height  string `json:"height,omitempty"`
								Round   int    `json:"round,omitempty"`
								BlockID struct {
									Hash          string `json:"hash,omitempty"`
									PartSetHeader struct {
										Total int    `json:"total,omitempty"`
										Hash  string `json:"hash,omitempty"`
									} `json:"part_set_header,omitempty"`
								} `json:"block_id,omitempty"`
								Signatures []struct {
									BlockIDFlag      string    `json:"block_id_flag,omitempty"`
									ValidatorAddress string    `json:"validator_address,omitempty"`
									Timestamp        time.Time `json:"timestamp,omitempty"`
									Signature        string    `json:"signature,omitempty"`
								} `json:"signatures,omitempty"`
							} `json:"commit,omitempty"`
						} `json:"signed_header,omitempty"`
						ValidatorSet struct {
							Validators []struct {
								Address string `json:"address,omitempty"`
								PubKey  struct {
									Ed25519   string `json:"ed25519,omitempty"`
									Secp256K1 string `json:"secp256k1,omitempty"`
								} `json:"pub_key,omitempty"`
								VotingPower      string `json:"voting_power,omitempty"`
								ProposerPriority string `json:"proposer_priority,omitempty"`
							} `json:"validators,omitempty"`
							Proposer struct {
								Address string `json:"address,omitempty"`
								PubKey  struct {
									Ed25519   string `json:"ed25519,omitempty"`
									Secp256K1 string `json:"secp256k1,omitempty"`
								} `json:"pub_key,omitempty"`
								VotingPower      string `json:"voting_power,omitempty"`
								ProposerPriority string `json:"proposer_priority,omitempty"`
							} `json:"proposer,omitempty"`
							TotalVotingPower string `json:"total_voting_power,omitempty"`
						} `json:"validator_set,omitempty"`
					} `json:"conflicting_block,omitempty"`
					CommonHeight        string `json:"common_height,omitempty"`
					ByzantineValidators []struct {
						Address string `json:"address,omitempty"`
						PubKey  struct {
							Ed25519   string `json:"ed25519,omitempty"`
							Secp256K1 string `json:"secp256k1,omitempty"`
						} `json:"pub_key,omitempty"`
						VotingPower      string `json:"voting_power,omitempty"`
						ProposerPriority string `json:"proposer_priority,omitempty"`
					} `json:"byzantine_validators,omitempty"`
					TotalVotingPower string    `json:"total_voting_power,omitempty"`
					Timestamp        time.Time `json:"timestamp,omitempty"`
				} `json:"light_client_attack_evidence,omitempty"`
			} `json:"evidence,omitempty"`
		} `json:"evidence,omitempty"`
		LastCommit struct {
			Height  string `json:"height,omitempty"`
			Round   int    `json:"round,omitempty"`
			BlockID struct {
				Hash          string `json:"hash,omitempty"`
				PartSetHeader struct {
					Total int    `json:"total,omitempty"`
					Hash  string `json:"hash,omitempty"`
				} `json:"part_set_header,omitempty"`
			} `json:"block_id,omitempty"`
			Signatures []struct {
				BlockIDFlag      string    `json:"block_id_flag,omitempty"`
				ValidatorAddress string    `json:"validator_address,omitempty"`
				Timestamp        time.Time `json:"timestamp,omitempty"`
				Signature        string    `json:"signature,omitempty"`
			} `json:"signatures,omitempty"`
		} `json:"last_commit,omitempty"`
	} `json:"block,omitempty"`
	SdkBlock struct {
		Header struct {
			Version struct {
				Block string `json:"block,omitempty"`
				App   string `json:"app,omitempty"`
			} `json:"version,omitempty"`
			ChainID     string    `json:"chain_id,omitempty"`
			Height      string    `json:"height,omitempty"`
			Time        time.Time `json:"time,omitempty"`
			LastBlockID struct {
				Hash          string `json:"hash,omitempty"`
				PartSetHeader struct {
					Total int    `json:"total,omitempty"`
					Hash  string `json:"hash,omitempty"`
				} `json:"part_set_header,omitempty"`
			} `json:"last_block_id,omitempty"`
			LastCommitHash     string `json:"last_commit_hash,omitempty"`
			DataHash           string `json:"data_hash,omitempty"`
			ValidatorsHash     string `json:"validators_hash,omitempty"`
			NextValidatorsHash string `json:"next_validators_hash,omitempty"`
			ConsensusHash      string `json:"consensus_hash,omitempty"`
			AppHash            string `json:"app_hash,omitempty"`
			LastResultsHash    string `json:"last_results_hash,omitempty"`
			EvidenceHash       string `json:"evidence_hash,omitempty"`
			ProposerAddress    string `json:"proposer_address,omitempty"`
		} `json:"header,omitempty"`
		Data struct {
			Txs []string `json:"txs,omitempty"`
		} `json:"data,omitempty"`
		Evidence struct {
			Evidence []struct {
				DuplicateVoteEvidence struct {
					VoteA struct {
						Type    string `json:"type,omitempty"`
						Height  string `json:"height,omitempty"`
						Round   int    `json:"round,omitempty"`
						BlockID struct {
							Hash          string `json:"hash,omitempty"`
							PartSetHeader struct {
								Total int    `json:"total,omitempty"`
								Hash  string `json:"hash,omitempty"`
							} `json:"part_set_header,omitempty"`
						} `json:"block_id,omitempty"`
						Timestamp        time.Time `json:"timestamp,omitempty"`
						ValidatorAddress string    `json:"validator_address,omitempty"`
						ValidatorIndex   int       `json:"validator_index,omitempty"`
						Signature        string    `json:"signature,omitempty"`
					} `json:"vote_a,omitempty"`
					VoteB struct {
						Type    string `json:"type,omitempty"`
						Height  string `json:"height,omitempty"`
						Round   int    `json:"round,omitempty"`
						BlockID struct {
							Hash          string `json:"hash,omitempty"`
							PartSetHeader struct {
								Total int    `json:"total,omitempty"`
								Hash  string `json:"hash,omitempty"`
							} `json:"part_set_header,omitempty"`
						} `json:"block_id,omitempty"`
						Timestamp        time.Time `json:"timestamp,omitempty"`
						ValidatorAddress string    `json:"validator_address,omitempty"`
						ValidatorIndex   int       `json:"validator_index,omitempty"`
						Signature        string    `json:"signature,omitempty"`
					} `json:"vote_b,omitempty"`
					TotalVotingPower string    `json:"total_voting_power,omitempty"`
					ValidatorPower   string    `json:"validator_power,omitempty"`
					Timestamp        time.Time `json:"timestamp,omitempty"`
				} `json:"duplicate_vote_evidence,omitempty"`
				LightClientAttackEvidence struct {
					ConflictingBlock struct {
						SignedHeader struct {
							Header struct {
								Version struct {
									Block string `json:"block,omitempty"`
									App   string `json:"app,omitempty"`
								} `json:"version,omitempty"`
								ChainID     string    `json:"chain_id,omitempty"`
								Height      string    `json:"height,omitempty"`
								Time        time.Time `json:"time,omitempty"`
								LastBlockID struct {
									Hash          string `json:"hash,omitempty"`
									PartSetHeader struct {
										Total int    `json:"total,omitempty"`
										Hash  string `json:"hash,omitempty"`
									} `json:"part_set_header,omitempty"`
								} `json:"last_block_id,omitempty"`
								LastCommitHash     string `json:"last_commit_hash,omitempty"`
								DataHash           string `json:"data_hash,omitempty"`
								ValidatorsHash     string `json:"validators_hash,omitempty"`
								NextValidatorsHash string `json:"next_validators_hash,omitempty"`
								ConsensusHash      string `json:"consensus_hash,omitempty"`
								AppHash            string `json:"app_hash,omitempty"`
								LastResultsHash    string `json:"last_results_hash,omitempty"`
								EvidenceHash       string `json:"evidence_hash,omitempty"`
								ProposerAddress    string `json:"proposer_address,omitempty"`
							} `json:"header,omitempty"`
							Commit struct {
								Height  string `json:"height,omitempty"`
								Round   int    `json:"round,omitempty"`
								BlockID struct {
									Hash          string `json:"hash,omitempty"`
									PartSetHeader struct {
										Total int    `json:"total,omitempty"`
										Hash  string `json:"hash,omitempty"`
									} `json:"part_set_header,omitempty"`
								} `json:"block_id,omitempty"`
								Signatures []struct {
									BlockIDFlag      string    `json:"block_id_flag,omitempty"`
									ValidatorAddress string    `json:"validator_address,omitempty"`
									Timestamp        time.Time `json:"timestamp,omitempty"`
									Signature        string    `json:"signature,omitempty"`
								} `json:"signatures,omitempty"`
							} `json:"commit,omitempty"`
						} `json:"signed_header,omitempty"`
						ValidatorSet struct {
							Validators []struct {
								Address string `json:"address,omitempty"`
								PubKey  struct {
									Ed25519   string `json:"ed25519,omitempty"`
									Secp256K1 string `json:"secp256k1,omitempty"`
								} `json:"pub_key,omitempty"`
								VotingPower      string `json:"voting_power,omitempty"`
								ProposerPriority string `json:"proposer_priority,omitempty"`
							} `json:"validators,omitempty"`
							Proposer struct {
								Address string `json:"address,omitempty"`
								PubKey  struct {
									Ed25519   string `json:"ed25519,omitempty"`
									Secp256K1 string `json:"secp256k1,omitempty"`
								} `json:"pub_key,omitempty"`
								VotingPower      string `json:"voting_power,omitempty"`
								ProposerPriority string `json:"proposer_priority,omitempty"`
							} `json:"proposer,omitempty"`
							TotalVotingPower string `json:"total_voting_power,omitempty"`
						} `json:"validator_set,omitempty"`
					} `json:"conflicting_block,omitempty"`
					CommonHeight        string `json:"common_height,omitempty"`
					ByzantineValidators []struct {
						Address string `json:"address,omitempty"`
						PubKey  struct {
							Ed25519   string `json:"ed25519,omitempty"`
							Secp256K1 string `json:"secp256k1,omitempty"`
						} `json:"pub_key,omitempty"`
						VotingPower      string `json:"voting_power,omitempty"`
						ProposerPriority string `json:"proposer_priority,omitempty"`
					} `json:"byzantine_validators,omitempty"`
					TotalVotingPower string    `json:"total_voting_power,omitempty"`
					Timestamp        time.Time `json:"timestamp,omitempty"`
				} `json:"light_client_attack_evidence,omitempty"`
			} `json:"evidence,omitempty"`
		} `json:"evidence,omitempty"`
		LastCommit struct {
			Height  string `json:"height,omitempty"`
			Round   int    `json:"round,omitempty"`
			BlockID struct {
				Hash          string `json:"hash,omitempty"`
				PartSetHeader struct {
					Total int    `json:"total,omitempty"`
					Hash  string `json:"hash,omitempty"`
				} `json:"part_set_header,omitempty"`
			} `json:"block_id,omitempty"`
			Signatures []struct {
				BlockIDFlag      string    `json:"block_id_flag,omitempty"`
				ValidatorAddress string    `json:"validator_address,omitempty"`
				Timestamp        time.Time `json:"timestamp,omitempty"`
				Signature        string    `json:"signature,omitempty"`
			} `json:"signatures,omitempty"`
		} `json:"last_commit,omitempty"`
	} `json:"sdk_block,omitempty"`
}
