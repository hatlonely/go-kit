// autogen by github.com/hatlonely/go-kit/astx/wrap.go. do not edit!
package wrap

import (
	"context"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"

	"github.com/hatlonely/go-kit/config"
	"github.com/hatlonely/go-kit/refx"
)

type KMSClientWrapper struct {
	obj     *kms.Client
	retry   *Retry
	options *WrapperOptions
}

func (w *KMSClientWrapper) Unwrap() *kms.Client {
	return w.obj
}

func (w *KMSClientWrapper) OnWrapperChange(opts ...refx.Option) config.OnChangeHandler {
	return func(cfg *config.Config) error {
		var options WrapperOptions
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.Wrap(err, "cfg.Unmarshal failed")
		}
		w.options = &options
		return nil
	}
}

func (w *KMSClientWrapper) OnRetryChange(opts ...refx.Option) config.OnChangeHandler {
	return func(cfg *config.Config) error {
		var options RetryOptions
		if err := cfg.Unmarshal(&options, opts...); err != nil {
			return errors.Wrap(err, "cfg.Unmarshal failed")
		}
		retry, err := NewRetryWithOptions(&options)
		if err != nil {
			return errors.Wrap(err, "NewRetryWithOptions failed")
		}
		w.retry = retry
		return nil
	}
}

func (w *KMSClientWrapper) AsymmetricDecrypt(ctx context.Context, request *kms.AsymmetricDecryptRequest) (*kms.AsymmetricDecryptResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.AsymmetricDecrypt")
		defer span.Finish()
	}

	var response *kms.AsymmetricDecryptResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.AsymmetricDecrypt(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) AsymmetricDecryptWithCallback(ctx context.Context, request *kms.AsymmetricDecryptRequest, callback func(response *kms.AsymmetricDecryptResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.AsymmetricDecryptWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.AsymmetricDecryptWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) AsymmetricDecryptWithChan(ctx context.Context, request *kms.AsymmetricDecryptRequest) (<-chan *kms.AsymmetricDecryptResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.AsymmetricDecryptWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.AsymmetricDecryptWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) AsymmetricEncrypt(ctx context.Context, request *kms.AsymmetricEncryptRequest) (*kms.AsymmetricEncryptResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.AsymmetricEncrypt")
		defer span.Finish()
	}

	var response *kms.AsymmetricEncryptResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.AsymmetricEncrypt(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) AsymmetricEncryptWithCallback(ctx context.Context, request *kms.AsymmetricEncryptRequest, callback func(response *kms.AsymmetricEncryptResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.AsymmetricEncryptWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.AsymmetricEncryptWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) AsymmetricEncryptWithChan(ctx context.Context, request *kms.AsymmetricEncryptRequest) (<-chan *kms.AsymmetricEncryptResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.AsymmetricEncryptWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.AsymmetricEncryptWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) AsymmetricSign(ctx context.Context, request *kms.AsymmetricSignRequest) (*kms.AsymmetricSignResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.AsymmetricSign")
		defer span.Finish()
	}

	var response *kms.AsymmetricSignResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.AsymmetricSign(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) AsymmetricSignWithCallback(ctx context.Context, request *kms.AsymmetricSignRequest, callback func(response *kms.AsymmetricSignResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.AsymmetricSignWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.AsymmetricSignWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) AsymmetricSignWithChan(ctx context.Context, request *kms.AsymmetricSignRequest) (<-chan *kms.AsymmetricSignResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.AsymmetricSignWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.AsymmetricSignWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) AsymmetricVerify(ctx context.Context, request *kms.AsymmetricVerifyRequest) (*kms.AsymmetricVerifyResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.AsymmetricVerify")
		defer span.Finish()
	}

	var response *kms.AsymmetricVerifyResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.AsymmetricVerify(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) AsymmetricVerifyWithCallback(ctx context.Context, request *kms.AsymmetricVerifyRequest, callback func(response *kms.AsymmetricVerifyResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.AsymmetricVerifyWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.AsymmetricVerifyWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) AsymmetricVerifyWithChan(ctx context.Context, request *kms.AsymmetricVerifyRequest) (<-chan *kms.AsymmetricVerifyResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.AsymmetricVerifyWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.AsymmetricVerifyWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) CancelKeyDeletion(ctx context.Context, request *kms.CancelKeyDeletionRequest) (*kms.CancelKeyDeletionResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.CancelKeyDeletion")
		defer span.Finish()
	}

	var response *kms.CancelKeyDeletionResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.CancelKeyDeletion(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) CancelKeyDeletionWithCallback(ctx context.Context, request *kms.CancelKeyDeletionRequest, callback func(response *kms.CancelKeyDeletionResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.CancelKeyDeletionWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.CancelKeyDeletionWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) CancelKeyDeletionWithChan(ctx context.Context, request *kms.CancelKeyDeletionRequest) (<-chan *kms.CancelKeyDeletionResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.CancelKeyDeletionWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.CancelKeyDeletionWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) CertificatePrivateKeyDecrypt(ctx context.Context, request *kms.CertificatePrivateKeyDecryptRequest) (*kms.CertificatePrivateKeyDecryptResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.CertificatePrivateKeyDecrypt")
		defer span.Finish()
	}

	var response *kms.CertificatePrivateKeyDecryptResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.CertificatePrivateKeyDecrypt(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) CertificatePrivateKeyDecryptWithCallback(ctx context.Context, request *kms.CertificatePrivateKeyDecryptRequest, callback func(response *kms.CertificatePrivateKeyDecryptResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.CertificatePrivateKeyDecryptWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.CertificatePrivateKeyDecryptWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) CertificatePrivateKeyDecryptWithChan(ctx context.Context, request *kms.CertificatePrivateKeyDecryptRequest) (<-chan *kms.CertificatePrivateKeyDecryptResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.CertificatePrivateKeyDecryptWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.CertificatePrivateKeyDecryptWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) CertificatePrivateKeySign(ctx context.Context, request *kms.CertificatePrivateKeySignRequest) (*kms.CertificatePrivateKeySignResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.CertificatePrivateKeySign")
		defer span.Finish()
	}

	var response *kms.CertificatePrivateKeySignResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.CertificatePrivateKeySign(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) CertificatePrivateKeySignWithCallback(ctx context.Context, request *kms.CertificatePrivateKeySignRequest, callback func(response *kms.CertificatePrivateKeySignResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.CertificatePrivateKeySignWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.CertificatePrivateKeySignWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) CertificatePrivateKeySignWithChan(ctx context.Context, request *kms.CertificatePrivateKeySignRequest) (<-chan *kms.CertificatePrivateKeySignResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.CertificatePrivateKeySignWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.CertificatePrivateKeySignWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) CertificatePublicKeyEncrypt(ctx context.Context, request *kms.CertificatePublicKeyEncryptRequest) (*kms.CertificatePublicKeyEncryptResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.CertificatePublicKeyEncrypt")
		defer span.Finish()
	}

	var response *kms.CertificatePublicKeyEncryptResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.CertificatePublicKeyEncrypt(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) CertificatePublicKeyEncryptWithCallback(ctx context.Context, request *kms.CertificatePublicKeyEncryptRequest, callback func(response *kms.CertificatePublicKeyEncryptResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.CertificatePublicKeyEncryptWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.CertificatePublicKeyEncryptWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) CertificatePublicKeyEncryptWithChan(ctx context.Context, request *kms.CertificatePublicKeyEncryptRequest) (<-chan *kms.CertificatePublicKeyEncryptResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.CertificatePublicKeyEncryptWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.CertificatePublicKeyEncryptWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) CertificatePublicKeyVerify(ctx context.Context, request *kms.CertificatePublicKeyVerifyRequest) (*kms.CertificatePublicKeyVerifyResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.CertificatePublicKeyVerify")
		defer span.Finish()
	}

	var response *kms.CertificatePublicKeyVerifyResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.CertificatePublicKeyVerify(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) CertificatePublicKeyVerifyWithCallback(ctx context.Context, request *kms.CertificatePublicKeyVerifyRequest, callback func(response *kms.CertificatePublicKeyVerifyResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.CertificatePublicKeyVerifyWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.CertificatePublicKeyVerifyWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) CertificatePublicKeyVerifyWithChan(ctx context.Context, request *kms.CertificatePublicKeyVerifyRequest) (<-chan *kms.CertificatePublicKeyVerifyResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.CertificatePublicKeyVerifyWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.CertificatePublicKeyVerifyWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) CreateAlias(ctx context.Context, request *kms.CreateAliasRequest) (*kms.CreateAliasResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.CreateAlias")
		defer span.Finish()
	}

	var response *kms.CreateAliasResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.CreateAlias(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) CreateAliasWithCallback(ctx context.Context, request *kms.CreateAliasRequest, callback func(response *kms.CreateAliasResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.CreateAliasWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.CreateAliasWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) CreateAliasWithChan(ctx context.Context, request *kms.CreateAliasRequest) (<-chan *kms.CreateAliasResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.CreateAliasWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.CreateAliasWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) CreateCertificate(ctx context.Context, request *kms.CreateCertificateRequest) (*kms.CreateCertificateResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.CreateCertificate")
		defer span.Finish()
	}

	var response *kms.CreateCertificateResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.CreateCertificate(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) CreateCertificateWithCallback(ctx context.Context, request *kms.CreateCertificateRequest, callback func(response *kms.CreateCertificateResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.CreateCertificateWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.CreateCertificateWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) CreateCertificateWithChan(ctx context.Context, request *kms.CreateCertificateRequest) (<-chan *kms.CreateCertificateResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.CreateCertificateWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.CreateCertificateWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) CreateKey(ctx context.Context, request *kms.CreateKeyRequest) (*kms.CreateKeyResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.CreateKey")
		defer span.Finish()
	}

	var response *kms.CreateKeyResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.CreateKey(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) CreateKeyVersion(ctx context.Context, request *kms.CreateKeyVersionRequest) (*kms.CreateKeyVersionResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.CreateKeyVersion")
		defer span.Finish()
	}

	var response *kms.CreateKeyVersionResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.CreateKeyVersion(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) CreateKeyVersionWithCallback(ctx context.Context, request *kms.CreateKeyVersionRequest, callback func(response *kms.CreateKeyVersionResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.CreateKeyVersionWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.CreateKeyVersionWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) CreateKeyVersionWithChan(ctx context.Context, request *kms.CreateKeyVersionRequest) (<-chan *kms.CreateKeyVersionResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.CreateKeyVersionWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.CreateKeyVersionWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) CreateKeyWithCallback(ctx context.Context, request *kms.CreateKeyRequest, callback func(response *kms.CreateKeyResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.CreateKeyWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.CreateKeyWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) CreateKeyWithChan(ctx context.Context, request *kms.CreateKeyRequest) (<-chan *kms.CreateKeyResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.CreateKeyWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.CreateKeyWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) CreateSecret(ctx context.Context, request *kms.CreateSecretRequest) (*kms.CreateSecretResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.CreateSecret")
		defer span.Finish()
	}

	var response *kms.CreateSecretResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.CreateSecret(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) CreateSecretWithCallback(ctx context.Context, request *kms.CreateSecretRequest, callback func(response *kms.CreateSecretResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.CreateSecretWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.CreateSecretWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) CreateSecretWithChan(ctx context.Context, request *kms.CreateSecretRequest) (<-chan *kms.CreateSecretResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.CreateSecretWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.CreateSecretWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) Decrypt(ctx context.Context, request *kms.DecryptRequest) (*kms.DecryptResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.Decrypt")
		defer span.Finish()
	}

	var response *kms.DecryptResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.Decrypt(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) DecryptWithCallback(ctx context.Context, request *kms.DecryptRequest, callback func(response *kms.DecryptResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.DecryptWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.DecryptWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) DecryptWithChan(ctx context.Context, request *kms.DecryptRequest) (<-chan *kms.DecryptResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.DecryptWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.DecryptWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) DeleteAlias(ctx context.Context, request *kms.DeleteAliasRequest) (*kms.DeleteAliasResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.DeleteAlias")
		defer span.Finish()
	}

	var response *kms.DeleteAliasResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.DeleteAlias(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) DeleteAliasWithCallback(ctx context.Context, request *kms.DeleteAliasRequest, callback func(response *kms.DeleteAliasResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.DeleteAliasWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.DeleteAliasWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) DeleteAliasWithChan(ctx context.Context, request *kms.DeleteAliasRequest) (<-chan *kms.DeleteAliasResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.DeleteAliasWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.DeleteAliasWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) DeleteCertificate(ctx context.Context, request *kms.DeleteCertificateRequest) (*kms.DeleteCertificateResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.DeleteCertificate")
		defer span.Finish()
	}

	var response *kms.DeleteCertificateResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.DeleteCertificate(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) DeleteCertificateWithCallback(ctx context.Context, request *kms.DeleteCertificateRequest, callback func(response *kms.DeleteCertificateResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.DeleteCertificateWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.DeleteCertificateWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) DeleteCertificateWithChan(ctx context.Context, request *kms.DeleteCertificateRequest) (<-chan *kms.DeleteCertificateResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.DeleteCertificateWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.DeleteCertificateWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) DeleteKeyMaterial(ctx context.Context, request *kms.DeleteKeyMaterialRequest) (*kms.DeleteKeyMaterialResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.DeleteKeyMaterial")
		defer span.Finish()
	}

	var response *kms.DeleteKeyMaterialResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.DeleteKeyMaterial(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) DeleteKeyMaterialWithCallback(ctx context.Context, request *kms.DeleteKeyMaterialRequest, callback func(response *kms.DeleteKeyMaterialResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.DeleteKeyMaterialWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.DeleteKeyMaterialWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) DeleteKeyMaterialWithChan(ctx context.Context, request *kms.DeleteKeyMaterialRequest) (<-chan *kms.DeleteKeyMaterialResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.DeleteKeyMaterialWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.DeleteKeyMaterialWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) DeleteSecret(ctx context.Context, request *kms.DeleteSecretRequest) (*kms.DeleteSecretResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.DeleteSecret")
		defer span.Finish()
	}

	var response *kms.DeleteSecretResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.DeleteSecret(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) DeleteSecretWithCallback(ctx context.Context, request *kms.DeleteSecretRequest, callback func(response *kms.DeleteSecretResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.DeleteSecretWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.DeleteSecretWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) DeleteSecretWithChan(ctx context.Context, request *kms.DeleteSecretRequest) (<-chan *kms.DeleteSecretResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.DeleteSecretWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.DeleteSecretWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) DescribeAccountKmsStatus(ctx context.Context, request *kms.DescribeAccountKmsStatusRequest) (*kms.DescribeAccountKmsStatusResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.DescribeAccountKmsStatus")
		defer span.Finish()
	}

	var response *kms.DescribeAccountKmsStatusResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.DescribeAccountKmsStatus(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) DescribeAccountKmsStatusWithCallback(ctx context.Context, request *kms.DescribeAccountKmsStatusRequest, callback func(response *kms.DescribeAccountKmsStatusResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.DescribeAccountKmsStatusWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.DescribeAccountKmsStatusWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) DescribeAccountKmsStatusWithChan(ctx context.Context, request *kms.DescribeAccountKmsStatusRequest) (<-chan *kms.DescribeAccountKmsStatusResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.DescribeAccountKmsStatusWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.DescribeAccountKmsStatusWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) DescribeCertificate(ctx context.Context, request *kms.DescribeCertificateRequest) (*kms.DescribeCertificateResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.DescribeCertificate")
		defer span.Finish()
	}

	var response *kms.DescribeCertificateResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.DescribeCertificate(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) DescribeCertificateWithCallback(ctx context.Context, request *kms.DescribeCertificateRequest, callback func(response *kms.DescribeCertificateResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.DescribeCertificateWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.DescribeCertificateWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) DescribeCertificateWithChan(ctx context.Context, request *kms.DescribeCertificateRequest) (<-chan *kms.DescribeCertificateResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.DescribeCertificateWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.DescribeCertificateWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) DescribeKey(ctx context.Context, request *kms.DescribeKeyRequest) (*kms.DescribeKeyResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.DescribeKey")
		defer span.Finish()
	}

	var response *kms.DescribeKeyResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.DescribeKey(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) DescribeKeyVersion(ctx context.Context, request *kms.DescribeKeyVersionRequest) (*kms.DescribeKeyVersionResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.DescribeKeyVersion")
		defer span.Finish()
	}

	var response *kms.DescribeKeyVersionResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.DescribeKeyVersion(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) DescribeKeyVersionWithCallback(ctx context.Context, request *kms.DescribeKeyVersionRequest, callback func(response *kms.DescribeKeyVersionResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.DescribeKeyVersionWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.DescribeKeyVersionWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) DescribeKeyVersionWithChan(ctx context.Context, request *kms.DescribeKeyVersionRequest) (<-chan *kms.DescribeKeyVersionResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.DescribeKeyVersionWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.DescribeKeyVersionWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) DescribeKeyWithCallback(ctx context.Context, request *kms.DescribeKeyRequest, callback func(response *kms.DescribeKeyResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.DescribeKeyWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.DescribeKeyWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) DescribeKeyWithChan(ctx context.Context, request *kms.DescribeKeyRequest) (<-chan *kms.DescribeKeyResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.DescribeKeyWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.DescribeKeyWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) DescribeRegions(ctx context.Context, request *kms.DescribeRegionsRequest) (*kms.DescribeRegionsResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.DescribeRegions")
		defer span.Finish()
	}

	var response *kms.DescribeRegionsResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.DescribeRegions(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) DescribeRegionsWithCallback(ctx context.Context, request *kms.DescribeRegionsRequest, callback func(response *kms.DescribeRegionsResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.DescribeRegionsWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.DescribeRegionsWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) DescribeRegionsWithChan(ctx context.Context, request *kms.DescribeRegionsRequest) (<-chan *kms.DescribeRegionsResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.DescribeRegionsWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.DescribeRegionsWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) DescribeSecret(ctx context.Context, request *kms.DescribeSecretRequest) (*kms.DescribeSecretResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.DescribeSecret")
		defer span.Finish()
	}

	var response *kms.DescribeSecretResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.DescribeSecret(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) DescribeSecretWithCallback(ctx context.Context, request *kms.DescribeSecretRequest, callback func(response *kms.DescribeSecretResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.DescribeSecretWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.DescribeSecretWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) DescribeSecretWithChan(ctx context.Context, request *kms.DescribeSecretRequest) (<-chan *kms.DescribeSecretResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.DescribeSecretWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.DescribeSecretWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) DescribeService(ctx context.Context, request *kms.DescribeServiceRequest) (*kms.DescribeServiceResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.DescribeService")
		defer span.Finish()
	}

	var response *kms.DescribeServiceResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.DescribeService(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) DescribeServiceWithCallback(ctx context.Context, request *kms.DescribeServiceRequest, callback func(response *kms.DescribeServiceResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.DescribeServiceWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.DescribeServiceWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) DescribeServiceWithChan(ctx context.Context, request *kms.DescribeServiceRequest) (<-chan *kms.DescribeServiceResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.DescribeServiceWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.DescribeServiceWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) DisableKey(ctx context.Context, request *kms.DisableKeyRequest) (*kms.DisableKeyResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.DisableKey")
		defer span.Finish()
	}

	var response *kms.DisableKeyResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.DisableKey(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) DisableKeyWithCallback(ctx context.Context, request *kms.DisableKeyRequest, callback func(response *kms.DisableKeyResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.DisableKeyWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.DisableKeyWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) DisableKeyWithChan(ctx context.Context, request *kms.DisableKeyRequest) (<-chan *kms.DisableKeyResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.DisableKeyWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.DisableKeyWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) EnableKey(ctx context.Context, request *kms.EnableKeyRequest) (*kms.EnableKeyResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.EnableKey")
		defer span.Finish()
	}

	var response *kms.EnableKeyResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.EnableKey(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) EnableKeyWithCallback(ctx context.Context, request *kms.EnableKeyRequest, callback func(response *kms.EnableKeyResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.EnableKeyWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.EnableKeyWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) EnableKeyWithChan(ctx context.Context, request *kms.EnableKeyRequest) (<-chan *kms.EnableKeyResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.EnableKeyWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.EnableKeyWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) Encrypt(ctx context.Context, request *kms.EncryptRequest) (*kms.EncryptResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.Encrypt")
		defer span.Finish()
	}

	var response *kms.EncryptResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.Encrypt(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) EncryptWithCallback(ctx context.Context, request *kms.EncryptRequest, callback func(response *kms.EncryptResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.EncryptWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.EncryptWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) EncryptWithChan(ctx context.Context, request *kms.EncryptRequest) (<-chan *kms.EncryptResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.EncryptWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.EncryptWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) ExportCertificate(ctx context.Context, request *kms.ExportCertificateRequest) (*kms.ExportCertificateResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ExportCertificate")
		defer span.Finish()
	}

	var response *kms.ExportCertificateResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.ExportCertificate(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) ExportCertificateWithCallback(ctx context.Context, request *kms.ExportCertificateRequest, callback func(response *kms.ExportCertificateResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ExportCertificateWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.ExportCertificateWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) ExportCertificateWithChan(ctx context.Context, request *kms.ExportCertificateRequest) (<-chan *kms.ExportCertificateResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ExportCertificateWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.ExportCertificateWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) ExportDataKey(ctx context.Context, request *kms.ExportDataKeyRequest) (*kms.ExportDataKeyResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ExportDataKey")
		defer span.Finish()
	}

	var response *kms.ExportDataKeyResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.ExportDataKey(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) ExportDataKeyWithCallback(ctx context.Context, request *kms.ExportDataKeyRequest, callback func(response *kms.ExportDataKeyResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ExportDataKeyWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.ExportDataKeyWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) ExportDataKeyWithChan(ctx context.Context, request *kms.ExportDataKeyRequest) (<-chan *kms.ExportDataKeyResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ExportDataKeyWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.ExportDataKeyWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) GenerateAndExportDataKey(ctx context.Context, request *kms.GenerateAndExportDataKeyRequest) (*kms.GenerateAndExportDataKeyResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.GenerateAndExportDataKey")
		defer span.Finish()
	}

	var response *kms.GenerateAndExportDataKeyResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.GenerateAndExportDataKey(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) GenerateAndExportDataKeyWithCallback(ctx context.Context, request *kms.GenerateAndExportDataKeyRequest, callback func(response *kms.GenerateAndExportDataKeyResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.GenerateAndExportDataKeyWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.GenerateAndExportDataKeyWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) GenerateAndExportDataKeyWithChan(ctx context.Context, request *kms.GenerateAndExportDataKeyRequest) (<-chan *kms.GenerateAndExportDataKeyResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.GenerateAndExportDataKeyWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.GenerateAndExportDataKeyWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) GenerateDataKey(ctx context.Context, request *kms.GenerateDataKeyRequest) (*kms.GenerateDataKeyResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.GenerateDataKey")
		defer span.Finish()
	}

	var response *kms.GenerateDataKeyResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.GenerateDataKey(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) GenerateDataKeyWithCallback(ctx context.Context, request *kms.GenerateDataKeyRequest, callback func(response *kms.GenerateDataKeyResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.GenerateDataKeyWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.GenerateDataKeyWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) GenerateDataKeyWithChan(ctx context.Context, request *kms.GenerateDataKeyRequest) (<-chan *kms.GenerateDataKeyResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.GenerateDataKeyWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.GenerateDataKeyWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) GenerateDataKeyWithoutPlaintext(ctx context.Context, request *kms.GenerateDataKeyWithoutPlaintextRequest) (*kms.GenerateDataKeyWithoutPlaintextResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.GenerateDataKeyWithoutPlaintext")
		defer span.Finish()
	}

	var response *kms.GenerateDataKeyWithoutPlaintextResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.GenerateDataKeyWithoutPlaintext(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) GenerateDataKeyWithoutPlaintextWithCallback(ctx context.Context, request *kms.GenerateDataKeyWithoutPlaintextRequest, callback func(response *kms.GenerateDataKeyWithoutPlaintextResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.GenerateDataKeyWithoutPlaintextWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.GenerateDataKeyWithoutPlaintextWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) GenerateDataKeyWithoutPlaintextWithChan(ctx context.Context, request *kms.GenerateDataKeyWithoutPlaintextRequest) (<-chan *kms.GenerateDataKeyWithoutPlaintextResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.GenerateDataKeyWithoutPlaintextWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.GenerateDataKeyWithoutPlaintextWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) GetCertificate(ctx context.Context, request *kms.GetCertificateRequest) (*kms.GetCertificateResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.GetCertificate")
		defer span.Finish()
	}

	var response *kms.GetCertificateResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.GetCertificate(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) GetCertificateWithCallback(ctx context.Context, request *kms.GetCertificateRequest, callback func(response *kms.GetCertificateResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.GetCertificateWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.GetCertificateWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) GetCertificateWithChan(ctx context.Context, request *kms.GetCertificateRequest) (<-chan *kms.GetCertificateResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.GetCertificateWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.GetCertificateWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) GetParametersForImport(ctx context.Context, request *kms.GetParametersForImportRequest) (*kms.GetParametersForImportResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.GetParametersForImport")
		defer span.Finish()
	}

	var response *kms.GetParametersForImportResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.GetParametersForImport(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) GetParametersForImportWithCallback(ctx context.Context, request *kms.GetParametersForImportRequest, callback func(response *kms.GetParametersForImportResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.GetParametersForImportWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.GetParametersForImportWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) GetParametersForImportWithChan(ctx context.Context, request *kms.GetParametersForImportRequest) (<-chan *kms.GetParametersForImportResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.GetParametersForImportWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.GetParametersForImportWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) GetPublicKey(ctx context.Context, request *kms.GetPublicKeyRequest) (*kms.GetPublicKeyResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.GetPublicKey")
		defer span.Finish()
	}

	var response *kms.GetPublicKeyResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.GetPublicKey(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) GetPublicKeyWithCallback(ctx context.Context, request *kms.GetPublicKeyRequest, callback func(response *kms.GetPublicKeyResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.GetPublicKeyWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.GetPublicKeyWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) GetPublicKeyWithChan(ctx context.Context, request *kms.GetPublicKeyRequest) (<-chan *kms.GetPublicKeyResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.GetPublicKeyWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.GetPublicKeyWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) GetRandomPassword(ctx context.Context, request *kms.GetRandomPasswordRequest) (*kms.GetRandomPasswordResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.GetRandomPassword")
		defer span.Finish()
	}

	var response *kms.GetRandomPasswordResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.GetRandomPassword(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) GetRandomPasswordWithCallback(ctx context.Context, request *kms.GetRandomPasswordRequest, callback func(response *kms.GetRandomPasswordResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.GetRandomPasswordWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.GetRandomPasswordWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) GetRandomPasswordWithChan(ctx context.Context, request *kms.GetRandomPasswordRequest) (<-chan *kms.GetRandomPasswordResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.GetRandomPasswordWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.GetRandomPasswordWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) GetSecretValue(ctx context.Context, request *kms.GetSecretValueRequest) (*kms.GetSecretValueResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.GetSecretValue")
		defer span.Finish()
	}

	var response *kms.GetSecretValueResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.GetSecretValue(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) GetSecretValueWithCallback(ctx context.Context, request *kms.GetSecretValueRequest, callback func(response *kms.GetSecretValueResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.GetSecretValueWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.GetSecretValueWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) GetSecretValueWithChan(ctx context.Context, request *kms.GetSecretValueRequest) (<-chan *kms.GetSecretValueResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.GetSecretValueWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.GetSecretValueWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) ImportCertificate(ctx context.Context, request *kms.ImportCertificateRequest) (*kms.ImportCertificateResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ImportCertificate")
		defer span.Finish()
	}

	var response *kms.ImportCertificateResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.ImportCertificate(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) ImportCertificateWithCallback(ctx context.Context, request *kms.ImportCertificateRequest, callback func(response *kms.ImportCertificateResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ImportCertificateWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.ImportCertificateWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) ImportCertificateWithChan(ctx context.Context, request *kms.ImportCertificateRequest) (<-chan *kms.ImportCertificateResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ImportCertificateWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.ImportCertificateWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) ImportEncryptionCertificate(ctx context.Context, request *kms.ImportEncryptionCertificateRequest) (*kms.ImportEncryptionCertificateResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ImportEncryptionCertificate")
		defer span.Finish()
	}

	var response *kms.ImportEncryptionCertificateResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.ImportEncryptionCertificate(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) ImportEncryptionCertificateWithCallback(ctx context.Context, request *kms.ImportEncryptionCertificateRequest, callback func(response *kms.ImportEncryptionCertificateResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ImportEncryptionCertificateWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.ImportEncryptionCertificateWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) ImportEncryptionCertificateWithChan(ctx context.Context, request *kms.ImportEncryptionCertificateRequest) (<-chan *kms.ImportEncryptionCertificateResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ImportEncryptionCertificateWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.ImportEncryptionCertificateWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) ImportKeyMaterial(ctx context.Context, request *kms.ImportKeyMaterialRequest) (*kms.ImportKeyMaterialResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ImportKeyMaterial")
		defer span.Finish()
	}

	var response *kms.ImportKeyMaterialResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.ImportKeyMaterial(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) ImportKeyMaterialWithCallback(ctx context.Context, request *kms.ImportKeyMaterialRequest, callback func(response *kms.ImportKeyMaterialResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ImportKeyMaterialWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.ImportKeyMaterialWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) ImportKeyMaterialWithChan(ctx context.Context, request *kms.ImportKeyMaterialRequest) (<-chan *kms.ImportKeyMaterialResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ImportKeyMaterialWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.ImportKeyMaterialWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) ListAliases(ctx context.Context, request *kms.ListAliasesRequest) (*kms.ListAliasesResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ListAliases")
		defer span.Finish()
	}

	var response *kms.ListAliasesResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.ListAliases(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) ListAliasesByKeyId(ctx context.Context, request *kms.ListAliasesByKeyIdRequest) (*kms.ListAliasesByKeyIdResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ListAliasesByKeyId")
		defer span.Finish()
	}

	var response *kms.ListAliasesByKeyIdResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.ListAliasesByKeyId(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) ListAliasesByKeyIdWithCallback(ctx context.Context, request *kms.ListAliasesByKeyIdRequest, callback func(response *kms.ListAliasesByKeyIdResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ListAliasesByKeyIdWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.ListAliasesByKeyIdWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) ListAliasesByKeyIdWithChan(ctx context.Context, request *kms.ListAliasesByKeyIdRequest) (<-chan *kms.ListAliasesByKeyIdResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ListAliasesByKeyIdWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.ListAliasesByKeyIdWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) ListAliasesWithCallback(ctx context.Context, request *kms.ListAliasesRequest, callback func(response *kms.ListAliasesResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ListAliasesWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.ListAliasesWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) ListAliasesWithChan(ctx context.Context, request *kms.ListAliasesRequest) (<-chan *kms.ListAliasesResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ListAliasesWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.ListAliasesWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) ListCertificates(ctx context.Context, request *kms.ListCertificatesRequest) (*kms.ListCertificatesResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ListCertificates")
		defer span.Finish()
	}

	var response *kms.ListCertificatesResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.ListCertificates(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) ListCertificatesWithCallback(ctx context.Context, request *kms.ListCertificatesRequest, callback func(response *kms.ListCertificatesResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ListCertificatesWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.ListCertificatesWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) ListCertificatesWithChan(ctx context.Context, request *kms.ListCertificatesRequest) (<-chan *kms.ListCertificatesResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ListCertificatesWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.ListCertificatesWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) ListKeyVersions(ctx context.Context, request *kms.ListKeyVersionsRequest) (*kms.ListKeyVersionsResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ListKeyVersions")
		defer span.Finish()
	}

	var response *kms.ListKeyVersionsResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.ListKeyVersions(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) ListKeyVersionsWithCallback(ctx context.Context, request *kms.ListKeyVersionsRequest, callback func(response *kms.ListKeyVersionsResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ListKeyVersionsWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.ListKeyVersionsWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) ListKeyVersionsWithChan(ctx context.Context, request *kms.ListKeyVersionsRequest) (<-chan *kms.ListKeyVersionsResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ListKeyVersionsWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.ListKeyVersionsWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) ListKeys(ctx context.Context, request *kms.ListKeysRequest) (*kms.ListKeysResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ListKeys")
		defer span.Finish()
	}

	var response *kms.ListKeysResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.ListKeys(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) ListKeysWithCallback(ctx context.Context, request *kms.ListKeysRequest, callback func(response *kms.ListKeysResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ListKeysWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.ListKeysWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) ListKeysWithChan(ctx context.Context, request *kms.ListKeysRequest) (<-chan *kms.ListKeysResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ListKeysWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.ListKeysWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) ListResourceTags(ctx context.Context, request *kms.ListResourceTagsRequest) (*kms.ListResourceTagsResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ListResourceTags")
		defer span.Finish()
	}

	var response *kms.ListResourceTagsResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.ListResourceTags(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) ListResourceTagsWithCallback(ctx context.Context, request *kms.ListResourceTagsRequest, callback func(response *kms.ListResourceTagsResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ListResourceTagsWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.ListResourceTagsWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) ListResourceTagsWithChan(ctx context.Context, request *kms.ListResourceTagsRequest) (<-chan *kms.ListResourceTagsResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ListResourceTagsWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.ListResourceTagsWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) ListSecretVersionIds(ctx context.Context, request *kms.ListSecretVersionIdsRequest) (*kms.ListSecretVersionIdsResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ListSecretVersionIds")
		defer span.Finish()
	}

	var response *kms.ListSecretVersionIdsResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.ListSecretVersionIds(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) ListSecretVersionIdsWithCallback(ctx context.Context, request *kms.ListSecretVersionIdsRequest, callback func(response *kms.ListSecretVersionIdsResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ListSecretVersionIdsWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.ListSecretVersionIdsWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) ListSecretVersionIdsWithChan(ctx context.Context, request *kms.ListSecretVersionIdsRequest) (<-chan *kms.ListSecretVersionIdsResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ListSecretVersionIdsWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.ListSecretVersionIdsWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) ListSecrets(ctx context.Context, request *kms.ListSecretsRequest) (*kms.ListSecretsResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ListSecrets")
		defer span.Finish()
	}

	var response *kms.ListSecretsResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.ListSecrets(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) ListSecretsWithCallback(ctx context.Context, request *kms.ListSecretsRequest, callback func(response *kms.ListSecretsResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ListSecretsWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.ListSecretsWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) ListSecretsWithChan(ctx context.Context, request *kms.ListSecretsRequest) (<-chan *kms.ListSecretsResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ListSecretsWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.ListSecretsWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) OpenKmsService(ctx context.Context, request *kms.OpenKmsServiceRequest) (*kms.OpenKmsServiceResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.OpenKmsService")
		defer span.Finish()
	}

	var response *kms.OpenKmsServiceResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.OpenKmsService(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) OpenKmsServiceWithCallback(ctx context.Context, request *kms.OpenKmsServiceRequest, callback func(response *kms.OpenKmsServiceResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.OpenKmsServiceWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.OpenKmsServiceWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) OpenKmsServiceWithChan(ctx context.Context, request *kms.OpenKmsServiceRequest) (<-chan *kms.OpenKmsServiceResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.OpenKmsServiceWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.OpenKmsServiceWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) PutSecretValue(ctx context.Context, request *kms.PutSecretValueRequest) (*kms.PutSecretValueResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.PutSecretValue")
		defer span.Finish()
	}

	var response *kms.PutSecretValueResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.PutSecretValue(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) PutSecretValueWithCallback(ctx context.Context, request *kms.PutSecretValueRequest, callback func(response *kms.PutSecretValueResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.PutSecretValueWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.PutSecretValueWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) PutSecretValueWithChan(ctx context.Context, request *kms.PutSecretValueRequest) (<-chan *kms.PutSecretValueResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.PutSecretValueWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.PutSecretValueWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) ReEncrypt(ctx context.Context, request *kms.ReEncryptRequest) (*kms.ReEncryptResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ReEncrypt")
		defer span.Finish()
	}

	var response *kms.ReEncryptResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.ReEncrypt(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) ReEncryptWithCallback(ctx context.Context, request *kms.ReEncryptRequest, callback func(response *kms.ReEncryptResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ReEncryptWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.ReEncryptWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) ReEncryptWithChan(ctx context.Context, request *kms.ReEncryptRequest) (<-chan *kms.ReEncryptResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ReEncryptWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.ReEncryptWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) RestoreSecret(ctx context.Context, request *kms.RestoreSecretRequest) (*kms.RestoreSecretResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.RestoreSecret")
		defer span.Finish()
	}

	var response *kms.RestoreSecretResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.RestoreSecret(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) RestoreSecretWithCallback(ctx context.Context, request *kms.RestoreSecretRequest, callback func(response *kms.RestoreSecretResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.RestoreSecretWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.RestoreSecretWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) RestoreSecretWithChan(ctx context.Context, request *kms.RestoreSecretRequest) (<-chan *kms.RestoreSecretResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.RestoreSecretWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.RestoreSecretWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) ScheduleKeyDeletion(ctx context.Context, request *kms.ScheduleKeyDeletionRequest) (*kms.ScheduleKeyDeletionResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ScheduleKeyDeletion")
		defer span.Finish()
	}

	var response *kms.ScheduleKeyDeletionResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.ScheduleKeyDeletion(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) ScheduleKeyDeletionWithCallback(ctx context.Context, request *kms.ScheduleKeyDeletionRequest, callback func(response *kms.ScheduleKeyDeletionResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ScheduleKeyDeletionWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.ScheduleKeyDeletionWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) ScheduleKeyDeletionWithChan(ctx context.Context, request *kms.ScheduleKeyDeletionRequest) (<-chan *kms.ScheduleKeyDeletionResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.ScheduleKeyDeletionWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.ScheduleKeyDeletionWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) TagResource(ctx context.Context, request *kms.TagResourceRequest) (*kms.TagResourceResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.TagResource")
		defer span.Finish()
	}

	var response *kms.TagResourceResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.TagResource(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) TagResourceWithCallback(ctx context.Context, request *kms.TagResourceRequest, callback func(response *kms.TagResourceResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.TagResourceWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.TagResourceWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) TagResourceWithChan(ctx context.Context, request *kms.TagResourceRequest) (<-chan *kms.TagResourceResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.TagResourceWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.TagResourceWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) UntagResource(ctx context.Context, request *kms.UntagResourceRequest) (*kms.UntagResourceResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.UntagResource")
		defer span.Finish()
	}

	var response *kms.UntagResourceResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.UntagResource(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) UntagResourceWithCallback(ctx context.Context, request *kms.UntagResourceRequest, callback func(response *kms.UntagResourceResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.UntagResourceWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.UntagResourceWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) UntagResourceWithChan(ctx context.Context, request *kms.UntagResourceRequest) (<-chan *kms.UntagResourceResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.UntagResourceWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.UntagResourceWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) UpdateAlias(ctx context.Context, request *kms.UpdateAliasRequest) (*kms.UpdateAliasResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.UpdateAlias")
		defer span.Finish()
	}

	var response *kms.UpdateAliasResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.UpdateAlias(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) UpdateAliasWithCallback(ctx context.Context, request *kms.UpdateAliasRequest, callback func(response *kms.UpdateAliasResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.UpdateAliasWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.UpdateAliasWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) UpdateAliasWithChan(ctx context.Context, request *kms.UpdateAliasRequest) (<-chan *kms.UpdateAliasResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.UpdateAliasWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.UpdateAliasWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) UpdateCertificateStatus(ctx context.Context, request *kms.UpdateCertificateStatusRequest) (*kms.UpdateCertificateStatusResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.UpdateCertificateStatus")
		defer span.Finish()
	}

	var response *kms.UpdateCertificateStatusResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.UpdateCertificateStatus(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) UpdateCertificateStatusWithCallback(ctx context.Context, request *kms.UpdateCertificateStatusRequest, callback func(response *kms.UpdateCertificateStatusResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.UpdateCertificateStatusWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.UpdateCertificateStatusWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) UpdateCertificateStatusWithChan(ctx context.Context, request *kms.UpdateCertificateStatusRequest) (<-chan *kms.UpdateCertificateStatusResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.UpdateCertificateStatusWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.UpdateCertificateStatusWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) UpdateKeyDescription(ctx context.Context, request *kms.UpdateKeyDescriptionRequest) (*kms.UpdateKeyDescriptionResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.UpdateKeyDescription")
		defer span.Finish()
	}

	var response *kms.UpdateKeyDescriptionResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.UpdateKeyDescription(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) UpdateKeyDescriptionWithCallback(ctx context.Context, request *kms.UpdateKeyDescriptionRequest, callback func(response *kms.UpdateKeyDescriptionResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.UpdateKeyDescriptionWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.UpdateKeyDescriptionWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) UpdateKeyDescriptionWithChan(ctx context.Context, request *kms.UpdateKeyDescriptionRequest) (<-chan *kms.UpdateKeyDescriptionResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.UpdateKeyDescriptionWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.UpdateKeyDescriptionWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) UpdateRotationPolicy(ctx context.Context, request *kms.UpdateRotationPolicyRequest) (*kms.UpdateRotationPolicyResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.UpdateRotationPolicy")
		defer span.Finish()
	}

	var response *kms.UpdateRotationPolicyResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.UpdateRotationPolicy(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) UpdateRotationPolicyWithCallback(ctx context.Context, request *kms.UpdateRotationPolicyRequest, callback func(response *kms.UpdateRotationPolicyResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.UpdateRotationPolicyWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.UpdateRotationPolicyWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) UpdateRotationPolicyWithChan(ctx context.Context, request *kms.UpdateRotationPolicyRequest) (<-chan *kms.UpdateRotationPolicyResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.UpdateRotationPolicyWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.UpdateRotationPolicyWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) UpdateSecret(ctx context.Context, request *kms.UpdateSecretRequest) (*kms.UpdateSecretResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.UpdateSecret")
		defer span.Finish()
	}

	var response *kms.UpdateSecretResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.UpdateSecret(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) UpdateSecretVersionStage(ctx context.Context, request *kms.UpdateSecretVersionStageRequest) (*kms.UpdateSecretVersionStageResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.UpdateSecretVersionStage")
		defer span.Finish()
	}

	var response *kms.UpdateSecretVersionStageResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.UpdateSecretVersionStage(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) UpdateSecretVersionStageWithCallback(ctx context.Context, request *kms.UpdateSecretVersionStageRequest, callback func(response *kms.UpdateSecretVersionStageResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.UpdateSecretVersionStageWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.UpdateSecretVersionStageWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) UpdateSecretVersionStageWithChan(ctx context.Context, request *kms.UpdateSecretVersionStageRequest) (<-chan *kms.UpdateSecretVersionStageResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.UpdateSecretVersionStageWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.UpdateSecretVersionStageWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) UpdateSecretWithCallback(ctx context.Context, request *kms.UpdateSecretRequest, callback func(response *kms.UpdateSecretResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.UpdateSecretWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.UpdateSecretWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) UpdateSecretWithChan(ctx context.Context, request *kms.UpdateSecretRequest) (<-chan *kms.UpdateSecretResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.UpdateSecretWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.UpdateSecretWithChan(request)
	return res0, res1
}

func (w *KMSClientWrapper) UploadCertificate(ctx context.Context, request *kms.UploadCertificateRequest) (*kms.UploadCertificateResponse, error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.UploadCertificate")
		defer span.Finish()
	}

	var response *kms.UploadCertificateResponse
	var err error
	err = w.retry.Do(func() error {
		response, err = w.obj.UploadCertificate(request)
		return err
	})
	return response, err
}

func (w *KMSClientWrapper) UploadCertificateWithCallback(ctx context.Context, request *kms.UploadCertificateRequest, callback func(response *kms.UploadCertificateResponse, err error)) <-chan int {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.UploadCertificateWithCallback")
		defer span.Finish()
	}

	res0 := w.obj.UploadCertificateWithCallback(request, callback)
	return res0
}

func (w *KMSClientWrapper) UploadCertificateWithChan(ctx context.Context, request *kms.UploadCertificateRequest) (<-chan *kms.UploadCertificateResponse, <-chan error) {
	if w.options.EnableTrace {
		span, _ := opentracing.StartSpanFromContext(ctx, "kms.Client.UploadCertificateWithChan")
		defer span.Finish()
	}

	res0, res1 := w.obj.UploadCertificateWithChan(request)
	return res0, res1
}
